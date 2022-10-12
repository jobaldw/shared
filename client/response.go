package client

import (
	"io"
	"net/http"
)

// The response from an HTTP request.
type Response struct {
	body io.ReadCloser

	// http status test. Example "200 OK"
	Status string

	// http status code. Example 200
	StatusCode int

	// the request that was received by a server or to be sent by a
	// client
	Request *http.Request
}

// GetBody
// returns the response body in its i/o reader stream.
func (r *Response) GetBody() io.Reader {
	return r.body
}

// GetBodyBytes
// reads the response body i/o and converts it to a byte array.
func (r *Response) GetBodyBytes() []byte {
	bodyBytes, _ := io.ReadAll(r.body)
	return bodyBytes
}

// GetBodyString
// parses the response body i/o interface into a string.
func (r *Response) GetBodyString() string {
	bodyBytes, _ := io.ReadAll(r.body)
	return string(bodyBytes)
}

// IsSuccessful
// checks if the response status code is 200 level.
func (r *Response) IsSuccessful() bool {
	return r.check(199, 300)
}

// IsClientError
// checks if the response status code is 400 level.
func (r *Response) IsClientError() bool {
	return r.check(399, 500)
}

// IsServerError
// checks if the response status code is 500 level.
func (r *Response) IsServerError() bool {
	return r.check(499, 600)
}

func (r *Response) check(min, max int) bool {
	return r.StatusCode > min && r.StatusCode < max
}
