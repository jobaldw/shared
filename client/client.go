package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client struct
type Client struct {
	URL    *url.URL
	Client *http.Client
}

// New client
func New(uri string, timeout time.Duration) Client {
	return Client{
		URL: &url.URL{
			Host: uri,
		},
		Client: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

// Get request
func (c *Client) Get(headers, parameters map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.URL.Host, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range parameters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return c.Client.Do(req)
}

// StringResp of response body
func StringResp(body io.ReadCloser) (string, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), err
}

// MarshalResp of response body
func MarshalResp(body io.ReadCloser, structure interface{}) error {
	return json.NewDecoder(body).Decode(structure)
}
