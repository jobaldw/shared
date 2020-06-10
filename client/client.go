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
	URL    *url.URL
	Client *http.Client
}

// New client
func New(uri string, timeout int) Client {
	return Client{
		URL: &url.URL{
			Host: uri,
		},
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// Post request
func (c *Client) Post(endpoint string, headers, parameters map[string]string, body io.Reader) (resp response.Response, err error) {
	url := fmt.Sprintf("%s%s", c.URL.Host, endpoint)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		err = fmt.Errorf("problem occured invoking %s: %s, %s", http.MethodPost, url, err)
		return
	}

	req = prepRequest(req, headers, parameters)
	return c.do(req)
}

// Get request
func (c *Client) Get(endpoint string, headers, parameters map[string]string) (resp response.Response, err error) {
	url := fmt.Sprintf("%s%s", c.URL.Host, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("problem occured invoking %s: %s, %s", http.MethodGet, url, err)
		return
	}

	req = prepRequest(req, headers, parameters)
	return c.do(req)
}

// Put request
func (c *Client) Put(endpoint string, headers, parameters map[string]string, body io.Reader) (resp response.Response, err error) {
	url := fmt.Sprintf("%s%s", c.URL.Host, endpoint)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		err = fmt.Errorf("problem occured invoking %s: %s, %s", http.MethodPut, url, err)
		return
	}

	req = prepRequest(req, headers, parameters)
	return c.do(req)
}

// Delete request
func (c *Client) Delete(endpoint string, headers, parameters map[string]string) (resp response.Response, err error) {
	url := fmt.Sprintf("%s%s", c.URL.Host, endpoint)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		err = fmt.Errorf("problem occured invoking %s: %s, %s", http.MethodDelete, url, err)
		return
	}

	req = prepRequest(req, headers, parameters)
	return c.do(req)
}

// helper functions
func (c *Client) do(req *http.Request) (resp response.Response, err error) {
	r, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()

	resp.Save(r)

	return
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
