package coinbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type OAuth struct {
	clientID     string
	clientSecret string

	// redirectURL is optional, if not set, the first redirect url
	// in your Coinbase redirect urls list will be used.
	redirectURL string
}

func OAuthService(clientID, clientSecret, redirectURL string) *OAuth {
	o := new(OAuth)

	if redirectURL != "" {
		o.redirectURL = redirectURL
	}

	o.clientID = clientID
	o.clientSecret = clientSecret

	return o
}

// GenerateAuthorizationURL constructs the URL clients use to authorize our service
// against Coinbase OAuth2 service.
//
// state: Optional An unguessable random string. It is used to protect against cross-site
// request forgery attacks.
// scope: Optional Comma separated list of permissions (scopes) your application requests
// access to.
func (o *OAuth) GenerateAuthorizationURL(state string, scope []string) (string, error) {
	oauthURL, err := url.Parse(CoinURL)
	if err != nil {
		return "", err
	}

	parameters := make(url.Values, 3)
	parameters.Set("response_type", "code")
	parameters.Set("client_id", o.clientID)

	if o.redirectURL != "" {
		parameters.Set("redirect_uri", o.redirectURL)
	}

	if state != "" {
		parameters.Set("state", state)
	}

	if scope != nil {
		parameters.Set("scope", strings.Join(scope, ","))
	}

	oauthURL.RawQuery = parameters.Encode()
	oauthURL.Path = AuthEndpoint

	return oauthURL.String(), nil
}

type OAuthTokens struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshToken uses a refresh token to vend a new set of tokens
func (o *OAuth) RefreshToken(token *OAuthTokens) (*OAuthTokens, error) {
	return o.getTokens(token.RefreshToken, "refresh_token")
}

// GetToken uses OAuth temporary code returned form authorization
// server to get new set of tokens.
func (o *OAuth) GetToken(code string) (*OAuthTokens, error) {
	return o.getTokens(code, "authorization_code")
}

func (o *OAuth) getTokens(code, grantType string) (*OAuthTokens, error) {
	v := new(url.Values)
	v.Set("grant_type", grantType)
	v.Set("client_id", o.clientID)
	v.Set("client_secret", o.clientSecret)

	if grantType == "authorization_code" {
		v.Set("code", code)
		v.Set("redirect_uri", o.redirectURL)
	} else {
		v.Set("refresh_token", code)
	}

	req, err := http.NewRequest("POST", CoinAPI +TokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	c := new(http.Client)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		apiErr := new(coinAPIErr)
		apiErr.StatusCode = resp.StatusCode

		if len(bytes) != 0 { // try to return resp body as coinAPIErr.
			err := json.Unmarshal(bytes, apiErr)
			if err != nil {
				return nil, err
			}
		}

		return nil, apiErr
	}

	t := new(OAuthTokens)
	return t, json.Unmarshal(bytes, t)
}

// AuthorizationCode represents the temporary code (which used to request
// a valid access and refresh token from authorization server) and a
// state string (optional) that must match the state string in the first
// redirect url send from your app/service to the client.
type AuthorizationCode struct {
	Code  string
	State string
}

// GetAuthorizationCode processes redirection URL returned from oauth2
// server to the client before hitting your application/service and
// returns AuthorizationCode.
func GetAuthorizationCode(r *http.Request) (*AuthorizationCode, error) {
	a := new(AuthorizationCode)
	values := r.URL.Query()

	if code := values.Get("code"); code == "" {
		return nil, errors.New("no or invalid oauth2 temporary code found in this url")
	} else {
		a.Code = code
	}

	if state := values.Get("state"); state != "" {
		a.State = state
	}

	return a, nil
}

func GenerateState() string {
	return "random string"
}

func RevokeToken(token string) error {
	v := new(url.Values)
	v.Set("token", token)

	req, err := http.NewRequest("POST", CoinAPI +RevokeEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		apiErr := new(coinAPIErr)
		apiErr.StatusCode = resp.StatusCode

		if len(bytes) != 0 { // try to return resp body as coinAPIErr.
			err := json.Unmarshal(bytes, apiErr)
			if err != nil {
				return err
			}
		}

		return apiErr
	}

	return nil
}
