package coinbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type authenticator interface {
	writeAuthHeaders(*http.Request) error
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type operations interface {
	get(string) ([]byte, error)
	post(string, io.Reader) ([]byte, error)
	put(string, io.Reader) ([]byte, error)
	del(string, io.Reader) ([]byte, error)
}

type crud struct {
	auth       authenticator
	httpClient httpClient
}

func (c *crud) get(endpoint string) ([]byte, error) {
	return c.request("GET", endpoint, strings.NewReader(""))
}

func (c *crud) post(endpoint string, body io.Reader) ([]byte, error) {
	return c.request("POST", endpoint, body)
}

func (c *crud) put(endpoint string, body io.Reader) ([]byte, error) {
	return c.request("PUT", endpoint, body)
}

func (c *crud) del(endpoint string, body io.Reader) ([]byte, error) {
	return c.request("DELETE", endpoint, body)
}

func (c *crud) request(method, endpoint string, body io.Reader) ([]byte, error) {
	req, err := c.createRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	return c.executeRequest(req)
}

func (c *crud) createRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	buf := new(bytes.Buffer)
	body = strings.NewReader("")

	tee := io.TeeReader(body, buf)

	detect, err := ioutil.ReadAll(tee)
	if err != nil {
		return nil, err
	}

	contentType := "application/json"
	if http.DetectContentType(detect) != contentType {
		contentType = "application/x-www-form-urlencoded"
	}

	req, err := http.NewRequest(method, CoinAPI+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "CoinbaseGo/v1")
	req.Header.Set("Content-Type", contentType)

	if c.auth != nil { // send un-authenticated requests
		if err := c.auth.writeAuthHeaders(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (c *crud) executeRequest(req *http.Request) ([]byte, error) {
	if c.httpClient == nil {
		return nil, errors.New("http client can't be nil")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		apiErr := new(coinAPIErr)
		apiErr.StatusCode = resp.StatusCode
		apiErr.ReqMethod = req.Method
		apiErr.ReqURL = req.URL.String()

		// try to return response body in the error
		_ = json.Unmarshal(b, apiErr)

		return nil, apiErr
	}

	return b, nil
}
