package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Client struct
type Client struct {
	URL    string
	Client *http.Client
}

// New client
func New(url string, timeout time.Duration) Client {
	return Client{
		URL: url,
		Client: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

// Get request
func (c *Client) Get(headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

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
