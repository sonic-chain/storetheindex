package httpclient

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	v0 "github.com/filecoin-project/storetheindex/api/v0"
)

// New creates a base URL and a new http.Client.  The default port is only used
// if baseURL does not contain a port.
func New(baseURL, resource string, options ...Option) (*url.URL, *http.Client, error) {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}
	if u.Scheme == "" {
		return nil, nil, errors.New("url missing scheme")
	}
	u.Path = resource

	var cfg clientConfig
	if err := cfg.apply(options...); err != nil {
		return nil, nil, err
	}

	if cfg.client != nil {
		return u, cfg.client, nil
	}
	cl := &http.Client{
		Timeout: cfg.timeout,
	}
	return u, cl, nil
}

func ReadErrorFrom(status int, r io.Reader) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return ReadError(status, body)
}

func ReadError(status int, body []byte) error {
	se := v0.NewError(errors.New(strings.TrimSpace(string(body))), status)
	return errors.New(se.Text())
}
