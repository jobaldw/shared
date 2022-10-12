/*
Package client implements a client with the ability to make HTTP requests
and handle their responses.

The package client is limited to the below request methods:
  - DELETE, GET, POST, PUT

All functions that require a context to be passed should be given one from
the service handler request to correctly handle cancellations.
*/
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jobaldw/shared/v2/config"
)

const packageKey = "client" // package logging key

// A simple client that will also handle http requests. Uses the
// "net/http "and "net/url" packages.
type Client struct {
	health string
	client *http.Client

	// key-value pairs in an HTTP header
	Headers http.Header

	// a parsed URL (technically a URI)
	URL *url.URL
}

// New
// creates a new client from the shared config.Client() struct.
func New(conf config.Client) (*Client, error) {
	url, err := url.Parse(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %s, could not create client", packageKey, err)
	}

	return &Client{
		URL:     url,
		health:  conf.Health,
		Headers: conf.Headers,
		client: &http.Client{
			Timeout: time.Duration(conf.Timeout) * time.Second,
		},
	}, nil
}

// IsReady
// uses the clients health endpoint to determine if its up and running.
func (c *Client) IsReady(ctx context.Context) (bool, error) {
	resp, err := c.GetWithContext(ctx, c.health, nil)
	if err != nil {
		return false, fmt.Errorf("%s: %s", packageKey, err)
	}
	return resp.IsSuccessful(), nil
}

// Get
// makes a GET method request to the client with a background context.
func (c *Client) Get(path string, params map[string][]string) (*Response, error) {
	return c.GetWithContext(context.Background(), path, params)
}

// GetWithContext
// makes a GET method request to the client with any passed in context.
func (c *Client) GetWithContext(ctx context.Context, path string, params map[string][]string) (*Response, error) {
	return c.do(ctx, http.MethodGet, path, params, nil)
}

// Post
// makes a POST method request to the client with a background context.
func (c *Client) Post(path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.PostWithContext(context.Background(), path, params, body)
}

// PostWithContext
// makes a POST method request to the client with any passed in context.
func (c *Client) PostWithContext(ctx context.Context, path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPost, path, params, body)
}

// Put
// makes a PUT method request to the client with a background context.
func (c *Client) Put(path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.PutWithContext(context.Background(), path, params, body)
}

// PutWithContext
// makes a PUT method request to the client with any passed in context.
func (c *Client) PutWithContext(ctx context.Context, path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPut, path, params, body)
}

// Delete
// makes a DELETE method request to the client with a background context.
func (c *Client) Delete(path string, params map[string][]string) (*Response, error) {
	return c.DeleteWithContext(context.Background(), path, params)
}

// DeleteWithContext
// makes a DELETE method request to the client with any passed in context.
func (c *Client) DeleteWithContext(ctx context.Context, path string, params map[string][]string) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, params, nil)
}

/********** helper functions **********/

// do
// builds and makes the client request using the "net/https" package with
// NewRequestWithContext().
func (c *Client) do(ctx context.Context, method, path string, params url.Values, payload interface{}) (*Response, error) {
	// build the request body
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", packageKey, err)
		}
		body = bytes.NewBuffer(b)
	}

	// build the request url
	c.URL.Path = path
	uri := c.URL.ResolveReference(c.URL)
	req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
	if err != nil {
		return nil, fmt.Errorf("%s: %s, could not build request", packageKey, err)
	}
	req.Header = c.Headers
	req.URL.RawQuery = params.Encode()

	// do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %s, could not make request", packageKey, err)
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		body:       resp.Body,
		Request:    resp.Request,
	}, nil
}
