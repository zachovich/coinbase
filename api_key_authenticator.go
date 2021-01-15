package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
)

type apiKeyAuthenticator struct {
	apiKey    string
	apiSecret string
}

func (a *apiKeyAuthenticator) writeAuthHeaders(req *http.Request) error {
	if a.apiKey == "" || a.apiSecret == "" {
		return errors.New("api key/secret can't be empty")
	}

	hash := hmac.New(sha256.New, []byte(a.apiSecret))
	var body []byte
	var err error

	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	defer req.Body.Close()

	now := getTimestamp(nil)
	message := now + req.Method + req.URL.RequestURI() + string(body)

	_, err = hash.Write([]byte(message))
	if err != nil {
		return err
	}

	sha := hex.EncodeToString(hash.Sum(nil))

	req.Header.Set("CB-ACCESS-KEY", a.apiKey)
	req.Header.Set("CB-ACCESS-SIGN", sha)
	req.Header.Set("CB-ACCESS-TIMESTAMP", now)

	return nil
}
