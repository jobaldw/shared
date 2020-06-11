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
	URL        *url.URL
	Parameters map[string]string
	Headers    map[string]string
	Client     *http.Client
}

// New client
func New(rel string, timeout int, headers, parameters map[string]string) (Client, error) {
	uri, err := url.Parse(rel)
	if err != nil {
		return Client{}, fmt.Errorf("could not create client, %s", err)
	}

	return Client{
		URL:        uri,
		Parameters: parameters,
		Headers:    headers,
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, err
}

// Post request
func (c *Client) Post(path string, body io.Reader) (response.Response, error) {
	return c.do(http.MethodPost, path, c.Headers, c.Parameters, body)
}

// Get request
func (c *Client) Get(path string) (response.Response, error) {
	return c.do(http.MethodGet, path, c.Headers, c.Parameters, nil)
}

// Put request
func (c *Client) Put(path string, body io.Reader) (response.Response, error) {
	return c.do(http.MethodPut, path, c.Headers, c.Parameters, body)
}

// Delete request
func (c *Client) Delete(path string) (response.Response, error) {
	return c.do(http.MethodDelete, path, c.Headers, c.Parameters, nil)
}

// AddParams to client
func (c *Client) AddParams(clean bool, params map[string]string) {
	if clean {
		c.Parameters = map[string]string{}
	}

	for k, v := range params {
		c.Parameters[k] = v
	}
}

// helper functions
func (c *Client) do(method, path string, headers, parameters map[string]string, body io.Reader) (response.Response, error) {
	resp := response.Response{}

	req, err := http.NewRequest(method, c.URL.String()+"/"+path, body)
	if err != nil {
		return resp, fmt.Errorf("problem occured building %s request: %s, %s", method, c.URL.String(), err)
	}

	req = prepRequest(req, headers, parameters)
	r, err := c.Client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("problem occured with request: %s", err)
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
