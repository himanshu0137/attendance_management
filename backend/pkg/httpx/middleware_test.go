package httpx_test

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"hunrx.com/attendance_management/pkg/httpx"
)

// TestKeycloakAuth builds a Keycloak issuer serving one RSA key in the JWKS
// and a signed RS256 token, then verifies the auth middleware end to end.
func TestKeycloakAuth(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	_, err = x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}

	// Serve the JWKS document.
	jwksDoc := `{"keys":[{"kid":"k0","n":"` + b64(priv.N) + `","e":"` + b64(uint64(65537)) + `"}]}`
	mux := http.NewServeMux()
	mux.HandleFunc("/protocol/openid-connect/certs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jwksDoc))
	})
	mux.HandleFunc("/protocol/openid-connect/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"access_token": "at", "expires_in": 3600})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	issuer := srv.URL

	// Build a signed token whose payload carries org_id.
	now := time.Now().Unix()
	payload := map[string]any{
		"sub":    "u1",
		"iss":    issuer,
		"org_id": "t1",
		"roles":  []string{"ADMIN", "VIEWER"},
		"exp":    now,
		"iat":    now,
	}
	claims, _ := json.Marshal(payload)
	// A real JWT header so jws.Verify can parse and verify it.
	header, _ := json.Marshal(map[string]any{"alg": "RS256", "typ": "JWT"})
	sigInput := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(claims)
	digest := sha256.Sum256([]byte(sigInput))
	sigHdr, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	sigB64 := base64.RawURLEncoding.EncodeToString(sigHdr)
	token := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(header),
		base64.RawURLEncoding.EncodeToString(claims),
		sigB64,
	}, ".")

	// A claim with an expired exp must be rejected.
	expired := strings.Replace(token, `"exp":`+strconv.FormatInt(now, 10), `"exp":`+strconv.FormatInt(now-1000000000, 10), 1)

	// A claim issued by another issuer must be rejected.
	other := strings.Replace(token, `"iss":`+strconv.Quote(issuer), `"iss":`+strconv.Quote("https://other.example"), 1)

	for _, tc := range []struct {
		name   string
		token  string
		status int
	}{
		{"ok", token, http.StatusOK},
		{"expired", expired, http.StatusUnauthorized},
		{"wrong_issuer", other, http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := gin.New()
			g.Use(httpx.KeycloakAuth(issuer))
			g.Any("/x", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"sub": c.GetString("sub")})
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			rec := httptest.NewRecorder()
			g.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Errorf("got %d, want %d", rec.Code, tc.status)
			}
		})
	}
}

func b64(i any) string {
	switch v := i.(type) {
	case *big.Int:
		return base64.RawStdEncoding.EncodeToString(v.Bytes())
	case uint64:
		return base64.RawStdEncoding.EncodeToString([]byte(strconv.FormatInt(int64(v), 10)))
	default:
		return base64.RawStdEncoding.EncodeToString([]byte(strconv.FormatInt(v.(int64), 10)))
	}
}
