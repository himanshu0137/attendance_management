package organization

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

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

// clientCredentialsGrant is an oauth2.TokenSource that POSTs the client
// credentials grant.
type clientCredentialsGrant struct {
	url    string
	clientID, secret string
}

func (g clientCredentialsGrant) Token() (*oauth2.Token, error) {
	form := url.Values{
		"grant_type":     {"client_credentials"},
		"client_id":      {g.clientID},
		"client_secret":  {g.secret},
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
		return nil, fmt.Errorf("keycloak: unexpected status %d from %s", resp.StatusCode, g.url)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parsed tokenResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("keycloak: invalid token response JSON: %w", err)
	}
	if parsed.AccessToken == "" {
		return nil, fmt.Errorf("keycloak: token response has no access_token")
	}
	return &oauth2.Token{
		AccessToken: parsed.AccessToken,
		TokenType:   parsed.TokenType,
		ExpiresIn:   parsed.ExpiresIn,
	}, nil
}

// ServiceAccountToken obtains a service-account token for the Keycloak realm.
func ServiceAccountToken(ctx context.Context, issuer, realm, clientID, clientSecret string) (*oauth2.Token, error) {
	endpoint := issuer + "/realms/" + url.PathEscape(realm) + "/protocol/openid-connect/token"
	return clientCredentialsGrant{url: endpoint, clientID: clientID, secret: clientSecret}.Token()
}

// CreateKeycloakUser creates a user in the Keycloak realm and returns the
// Keycloak sub (user id) for org.employees.
func CreateKeycloakUser(ctx context.Context, issuer, realm, clientID, clientSecret, username, role string) (string, error) {
	st, err := ServiceAccountToken(ctx, issuer, realm, clientID, clientSecret)
	if err != nil {
		return "", err
	}

	form := url.Values{
		"username": {username},
		"enabled":  {"true"},
		"role":     {role},
	}
	req, err := http.NewRequest(http.MethodPost,
		issuer+"/admin/realms/"+url.PathEscape(realm)+"/users", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Authorization", "Bearer "+st.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("keycloak: user creation failed with status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var user struct {
		ID string `json:"id"`
		RealmAccess struct {
			Roles []string `json:"roles"`
		} `json:"realm_access"`
	}
	if err := json.Unmarshal(body, &user); err != nil {
		return "", fmt.Errorf("keycloak: invalid user response: %w", err)
	}
	if user.ID == "" {
		return "", fmt.Errorf("keycloak: user created without an id")
	}
	return user.ID, nil
}
