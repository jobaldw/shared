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

// package logging key
const packageKey = "client"

// A simple client that will also handle http requests. Uses the net/http and net/url packages.
type Client struct {
	// health path used for client health and ready checks
	health string

	// The net/http client struct.
	client *http.Client

	// Key-value pairs in an HTTP header.
	Headers http.Header

	// A parsed URL (technically a URI)
	URL *url.URL
}

// New
//
//	Creates a new client from the shared config.Client() struct.
//	* @param conf: client configurations
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
//
//	Uses the clients health endpoint to determine if its up and running.
//	* @param path: client path that request will send to
func (c *Client) IsReady(ctx context.Context) (bool, error) {
	resp, err := c.GetWithContext(ctx, c.health, nil)
	if err != nil {
		return false, fmt.Errorf("%s: %s", packageKey, err)
	}
	return resp.IsSuccessful(), nil
}

// Get
//
//	Makes a GET request to the client with a background context.Context().
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
func (c *Client) Get(path string, params map[string][]string) (*Response, error) {
	return c.GetWithContext(context.Background(), path, params)
}

// GetWithContext
//
//	Makes a GET request to the client with any passed in context.Context().
//	* @param ctx: context used to handle any cancellations
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
func (c *Client) GetWithContext(ctx context.Context, path string, params map[string][]string) (*Response, error) {
	return c.do(ctx, http.MethodGet, path, params, nil)
}

// Post
//
//	Makes a POST request to the client with any passed in context.Context().
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
//	* @param body: request body
func (c *Client) Post(path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.PostWithContext(context.Background(), path, params, body)
}

// PostWithContext
//
//	Makes a POST request to the client with any passed in context.Context().
//	* @param ctx: context used to handle any cancellations
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
//	* @param body: request body
func (c *Client) PostWithContext(ctx context.Context, path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPost, path, params, body)
}

// Put
//
//	Makes a PUT request to the client with any passed in context.Context().
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
//	* @param body: request body
func (c *Client) Put(path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.PutWithContext(context.Background(), path, params, body)
}

// PutWithContext
//
//	Makes a PUT request to the client with any passed in context.Context().
//	* @param ctx: context used to handle any cancellations
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
//	* @param body: request body
func (c *Client) PutWithContext(ctx context.Context, path string, params map[string][]string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPut, path, params, body)
}

// Delete
//
//	Makes a DELETE request to the client with a background context.Context().
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
func (c *Client) Delete(path string, params map[string][]string) (*Response, error) {
	return c.DeleteWithContext(context.Background(), path, params)
}

// DeleteWithContext
//
//	Makes a DELETE request to the client with any passed in context.Context().
//	* @param ctx: context used to handle any cancellations
//	* @param path: client path that request will send to
//	* @param parameters: query parameters
func (c *Client) DeleteWithContext(ctx context.Context, path string, params map[string][]string) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, params, nil)
}

/********** Helper functions **********/

// do
//
// Builds and makes the client request using the net/https package using NewRequestWithContext(). If one of the
// functions without context parameters are used (Get, Put, Post, Delete) then the context will be set to the
// background, else it will use whatever is passed in.
//   - @param ctx: context used to handle any cancellations
//   - @param method: HTTP request method
//   - @param path: client path that request will send to
//   - @param parameters: query parameters
//   - @param body: request payload to send
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
