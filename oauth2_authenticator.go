package coinbase

import (
	"errors"
	"net/http"
)

type oauthAuthenticator struct {
	token string
}

func (o *oauthAuthenticator) writeAuthHeaders(req *http.Request) error {
	if o.token == "" {
		return errors.New("authentication token can't be empty")
	}

	req.Header.Set("Authorization", "Bearer " + o.token)

	return nil
}