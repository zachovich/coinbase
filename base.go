package coinbase

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	CoinURL        = "https://coinbase.com"
	CoinAPI        = "https://api.coinbase.com"
	AuthEndpoint   = "/oauth/authorize"
	TokenEndpoint  = "/oauth/token"
	RevokeEndpoint = "/oauth/revoke"
)

func getTimestamp(c *http.Client) string {
	if c == nil {
		return strconv.FormatInt(time.Now().UTC().Unix(), 10)
	}

	return ""
}

type coinAPIErr struct {
	StatusCode     int    `json:"-"`
	ReqMethod      string `json:"-"`
	ReqURL         string `json:"-"`
	Err            string `json:"error"`
	ErrDescription string `json:"error_description"`
	Warning        string `json:"warning,omitempty"`
}

func (e *coinAPIErr) Error() string {
	m := fmt.Sprintf("status code: %d request method: %s request url: %s", e.StatusCode, e.ReqMethod, e.ReqURL)

	if e.Err != "" {
		m += fmt.Sprintf(" err: %s description: %s", e.Err, e.ErrDescription)
	}

	if e.Warning != "" {
		m += fmt.Sprintf(" warning: %s", e.Warning)
	}

	return m
}
