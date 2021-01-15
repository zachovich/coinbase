package coinbase

import (
	"encoding/json"
	"net/http"
)

type client struct {
	crud operations
}

func NewClient() (*client, error) {
	httpClient := new(http.Client)
	return &client{&crud{auth: nil, httpClient: httpClient}}, nil
}

func NewAPIKeyClient(key, secret string) (*client, error) {
	auth := &apiKeyAuthenticator{
		apiKey:    key,
		apiSecret: secret,
	}

	httpClient := new(http.Client)

	return &client{&crud{auth: auth, httpClient: httpClient}}, nil
}

func NewOAuthClient() (*client, error) {
	return nil, nil
}

func (c *client) ShowCurrentUser() (*user, error) {
	b, err := c.crud.get("/v2/user")
	if err != nil {
		return nil, err
	}

	u := new(user)

	err = json.Unmarshal(b, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
