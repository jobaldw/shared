package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jobaldw/shared/config"
)

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
// 	Creates a new client from the shared config.Client() struct.
// 	* @param conf: client configurations
func New(conf config.Client) (*Client, error) {
	url, err := url.Parse(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("%s, could not create client", err)
	}

	return &Client{
		URL:     url,
		health:  conf.Health,
		Headers: conf.Headers,
		client: &http.Client{
			Timeout: time.Duration(conf.Timeout) * time.Second,
		},
	}, err
}

// Get
//	Makes a GET request to the client with a background context.Context().
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
func (c *Client) Get(path string, params map[string][]string) (*Response, error) {
	return c.do(context.Background(), http.MethodGet, path, params, nil)
}

// GetWithContext
//	Makes a GET request to the client with any passed in context.Context().
// 	* @param ctx: context used to handle any cancellations
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
func (c *Client) GetWithContext(ctx context.Context, path string, params map[string][]string) (*Response, error) {
	return c.do(ctx, http.MethodGet, path, params, nil)
}

// Post
//	Makes a POST request to the client with any passed in context.Context().
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
// 	* @param body: request body
func (c *Client) Post(path string, params map[string][]string, body io.Reader) (*Response, error) {
	return c.do(context.Background(), http.MethodPost, path, params, body)
}

// PostWithContext
//	Makes a POST request to the client with any passed in context.Context().
// 	* @param ctx: context used to handle any cancellations
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
// 	* @param body: request body
func (c *Client) PostWithContext(ctx context.Context, path string, params map[string][]string, body io.Reader) (*Response, error) {
	return c.do(ctx, http.MethodPost, path, params, body)
}

// Put
//	Makes a PUT request to the client with any passed in context.Context().
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
// 	* @param body: request body
func (c *Client) Put(path string, params map[string][]string, body io.Reader) (*Response, error) {
	return c.do(context.Background(), http.MethodPut, path, params, body)
}

// PutWithContext
//	Makes a PUT request to the client with any passed in context.Context().
// 	* @param ctx: context used to handle any cancellations
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
// 	* @param body: request body
func (c *Client) PutWithContext(ctx context.Context, path string, params map[string][]string, body io.Reader) (*Response, error) {
	return c.do(ctx, http.MethodPut, path, params, body)
}

// Delete
//	Makes a DELETE request to the client with a background context.Context().
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
func (c *Client) Delete(path string, params map[string][]string) (*Response, error) {
	return c.do(context.Background(), http.MethodDelete, path, params, nil)
}

// DeleteWithContext
//	Makes a DELETE request to the client with any passed in context.Context().
// 	* @param ctx: context used to handle any cancellations
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
func (c *Client) DeleteWithContext(ctx context.Context, path string, params map[string][]string) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, params, nil)
}

/********** Helper functions **********/

// do
//
// Builds and makes the client request using the net/https package using NewRequestWithContext(). If one of the
// functions without context parameters are used (Get, Put, Post, Delete) then the context will be set to the
// background, else it will use whatever is passed in.
// 	* @param ctx: context used to handle any cancellations
// 	* @param method: HTTP request method
// 	* @param path: client path that request will send to
// 	* @param parameters: query parameters
// 	* @param body: request payload to send
func (c *Client) do(ctx context.Context, method, path string, params url.Values, body io.Reader) (*Response, error) {
	// build the request
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.URL.String(), path), body)
	if err != nil {
		return nil, fmt.Errorf(`%s, could not build %s request "%s"`, err, method, req.RequestURI)
	}
	req.Header = c.Headers
	req.URL.RawQuery = params.Encode()

	// do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(`%s, could not make %s request "%s"`, err, method, req.RequestURI)
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		body:       resp.Body,
		Request:    resp.Request,
	}, nil
}
