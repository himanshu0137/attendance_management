package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// tokenResponse is the JSON body of a Keycloak token request.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
}

// tokenGrant is an oauth2.TokenSource that POSTs a form-encoded token request.
type tokenGrant struct {
	url    string
	clientID string
	secret  string
	form    url.Values
}

func (g tokenGrant) Token() (*oauth2.Token, error) {
	form := g.form
	if form == nil {
		form = url.Values{}
	}
	req, err := http.NewRequest(http.MethodPost, g.url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(g.clientID, g.secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth: unexpected status %d from %s", resp.StatusCode, g.url)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parsed tokenResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("auth: invalid token response JSON: %w", err)
	}
	if parsed.AccessToken == "" {
		return nil, fmt.Errorf("auth: token response has no access_token")
	}
	return &oauth2.Token{
		AccessToken: parsed.AccessToken,
		TokenType:   parsed.TokenType,
		ExpiresIn:   parsed.ExpiresIn,
	}, nil
}

// authEndpoint is the Keycloak token endpoint used for the direct access grant.
var authEndpoint = func(ctx context.Context, issuer string) string {
	return issuer + "/protocol/openid-connect/token"
}

// Login exchanges username/password for an access token (direct access grant).
func Login(c *gin.Context, issuer, clientID, clientSecret string) (*oauth2.Token, error) {
	form := url.Values{
		"grant_type":   {"authorization_code"},
		"client_id":    {clientID},
		"redirect_uri": {issuer + "/protocol/openid-connect/cbs"},
		"code":         {c.PostForm("username")},
		"client_secret": {clientSecret},
	}
	return tokenGrant{url: authEndpoint(c.Request.Context(), issuer), clientID: clientID, secret: clientSecret, form: form}.Token()
}

// Refresh exchanges a refresh token for a new token pair.
func Refresh(c *gin.Context, issuer, clientID, clientSecret string) (*oauth2.Token, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.PostForm("refresh_token")},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}
	return tokenGrant{url: authEndpoint(c.Request.Context(), issuer), clientID: clientID, secret: clientSecret, form: form}.Token()
}

// Logout revokes the refresh token against the Keycloak revoke endpoint.
func Logout(ctx context.Context, issuer, clientID, clientSecret, refreshToken string) error {
	form := url.Values{"token": {refreshToken}}
	req, err := http.NewRequest(http.MethodPost, issuer+"/protocol/openid-connect/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(clientID, clientSecret)
	_, err = http.DefaultClient.Do(req)
	return err
}
