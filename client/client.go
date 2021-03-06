package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jobaldw/shared/client/response"
)

// Client struct
type Client struct {
	Health  string
	Headers map[string]string

	URL    *url.URL
	Client *http.Client
}

// New client
func New(rel, health string, timeout int, headers map[string]string) (Client, error) {
	uri, err := url.Parse(rel)
	if err != nil {
		return Client{}, fmt.Errorf("could not create client, %s", err)
	}

	return Client{
		URL:     uri,
		Health:  health,
		Headers: headers,
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, err
}

// Post request
func (c *Client) Post(path string, params map[string]string, body io.Reader) (response.Response, error) {
	return c.do(http.MethodPost, path, params, body)
}

// Get request
func (c *Client) Get(path string, params map[string]string) (response.Response, error) {
	return c.do(http.MethodGet, path, params, nil)
}

// Put request
func (c *Client) Put(path string, params map[string]string, body io.Reader) (response.Response, error) {
	return c.do(http.MethodPut, path, params, body)
}

// Delete request
func (c *Client) Delete(path string, params map[string]string) (response.Response, error) {
	return c.do(http.MethodDelete, path, params, nil)
}

// helper functions
func (c *Client) do(method, path string, parameters map[string]string, body io.Reader) (response.Response, error) {
	resp := response.Response{}

	req, err := http.NewRequest(method, c.URL.String()+"/"+path, body)
	if err != nil {
		return resp, fmt.Errorf("problem occurred building %s request: %s, %s", method, c.URL.String(), err)
	}

	req = prepRequest(req, c.Headers, parameters)
	r, err := c.Client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("problem occurred with request: %s", err)
	}

	resp.Save(r)

	return resp, nil
}

func prepRequest(req *http.Request, headers, parameters map[string]string) *http.Request {
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range parameters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req
}
