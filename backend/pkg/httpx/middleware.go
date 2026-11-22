package httpx

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/jws"
)

// KeycloakConfig is the Keycloak issuer identity used for local JWKS
// verification and the admin client credentials grant.
type KeycloakConfig struct {
	Issuer      string
	ClientID    string
	ClientSecret string
}

// KeycloakClaims is the set of claims injected into a request after local
// JWKS verification. org_id is the tenant identifier.
type KeycloakClaims struct {
	Sub    string            `json:"sub"`
	OrgID  string            `json:"org_id"`
	Roles  []string          `json:"roles"`
	Extra  map[string]any    `json:"-"`
	Raw    json.RawMessage   `json:"-"`
	Expiry time.Time         `json:"-"`
	Issuer string            `json:"-"`
}

// jwksKey is one RSA public key from the JWKS response.
type jwksKey struct {
	kid string
	pub *rsa.PublicKey
}

// jwks is a small JWKS cache keyed by kid.
type jwks struct {
	mu   sync.Mutex
	keys map[string]jwksKey
}

var jwksCache *jwks

// fetchKeySet GETs the JWKS document from the Keycloak issuer and caches it.
func fetchKeySet(ctx context.Context, issuer string) (*jwks, error) {
	resp, err := http.Get(issuer + "/protocol/openid-connect/certs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks: unexpected status %d from %s", resp.StatusCode, issuer+"/protocol/openid-connect/certs")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Keys []struct {
			KID string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("jwks: invalid document: %w", err)
	}
	j := &jwks{keys: map[string]jwksKey{}}
	for _, k := range doc.Keys {
		if k.KID == "" {
			continue
		}
		n, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			return nil, fmt.Errorf("jwks: invalid modulus for %q: %w", k.KID, err)
		}
		e, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			return nil, fmt.Errorf("jwks: invalid exponent for %q: %w", k.KID, err)
		}
		pub := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(e[0])}
		j.keys[k.KID] = jwksKey{kid: k.KID, pub: pub}
	}
	if len(j.keys) == 0 {
		return nil, fmt.Errorf("jwks: no RSA public keys found")
	}
	return j, nil
}

// KeycloakAuth guards every /api/<module>/* route. It verifies the bearer
// token against the Keycloak JWKS, injects sub/org_id/roles into the Gin
// context and enforces RBAC via RequireRole.
func KeycloakAuth(issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			Error(c, http.StatusUnauthorized, "missing_authorization", "Authorization header must start with \"Bearer \"")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := verifyToken(context.Background(), issuer, token)
		if err != nil {
			println("DEBUG verifyToken err: " + err.Error())
			Error(c, http.StatusUnauthorized, "invalid_token", err.Error())
			c.Abort()
			return
		}
		c.Set(KeycloakClaimsKey, claims)
		c.Set("sub", claims.Sub)
		c.Set("org_id", claims.OrgID)
		c.Set("roles", claims.Roles)
	}
}

// RequireRole allows only the given roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		found := false
		for _, r := range c.GetStringSlice("roles") {
			for _, want := range roles {
				if r == want {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			Error(c, http.StatusForbidden, "role_required", "insufficient roles")
			c.Abort()
		}
	}
}

// TenantID extracts the tenant identifier from the authenticated claims.
func TenantID(c *gin.Context) (string, bool) {
	v, ok := c.Get("org_id")
	return v.(string), ok
}

// KeycloakClaimsKey is the Gin context key for the verified claims.
var KeycloakClaimsKey = "keycloak_claims"

// verifyToken decodes and RS256-verifies a bearer token against the JWKS.
func verifyToken(ctx context.Context, issuer, token string) (KeycloakClaims, error) {
	claims, err := parseToken(context.Background(), issuer, token)
	if err != nil {
		return KeycloakClaims{}, err
	}
	if claims.Expiry.IsZero() || claims.Expiry.Before(time.Now()) {
		return KeycloakClaims{}, errors.New("token expired")
	}
	if claims.Issuer != issuer {
		return KeycloakClaims{}, fmt.Errorf("token issued by %q, expected %q", claims.Issuer, issuer)
	}
	return claims, nil
}

func parseToken(ctx context.Context, issuer, token string) (KeycloakClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return KeycloakClaims{}, fmt.Errorf("invalid token, expected 3 dot-separated parts")
	}
	var claimsJSON string
	for i, p := range parts {
		dec, err := base64.RawURLEncoding.DecodeString(p)
		if err != nil {
			return KeycloakClaims{}, fmt.Errorf("invalid base64url part %d: %w", i, err)
		}
		if i == 1 {
			claimsJSON = string(dec)
		}
	}
	var parsed struct {
		Sub   string   `json:"sub"`
		OrgID string   `json:"org_id"`
		Roles []string `json:"roles"`
		Exp   int64    `json:"exp"`
		Iss   string   `json:"iss"`
	}
	if err := json.Unmarshal([]byte(claimsJSON), &parsed); err != nil {
		return KeycloakClaims{}, fmt.Errorf("invalid claims JSON: %w", err)
	}
	if parsed.Sub == "" || parsed.OrgID == "" || parsed.Exp == 0 || parsed.Iss == "" {
		return KeycloakClaims{}, fmt.Errorf("claims missing sub, org_id, exp and iss")
	}
	exp := time.Unix(parsed.Exp, 0)
	claims := KeycloakClaims{
		Sub:    parsed.Sub,
		OrgID:  parsed.OrgID,
		Roles:  parsed.Roles,
		Expiry: exp,
		Issuer: parsed.Iss,
		Raw:    json.RawMessage(claimsJSON),
	}

	// Verify the signature locally against the cached JWKS key.
	keySet, err := fetchKeySet(ctx, issuer)
	if err != nil {
		return KeycloakClaims{}, err
	}
	var pub *rsa.PublicKey
	for _, k := range keySet.keys {
		if err := jws.Verify(token, k.pub); err == nil {
			pub = k.pub
			break
		}
	}
	if pub == nil {
		return KeycloakClaims{}, fmt.Errorf("no JWKS key signed this token")
	}
	return claims, nil
}
