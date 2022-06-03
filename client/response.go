package client

import (
	"io"
	"io/ioutil"
	"net/http"
)

// The response from an HTTP request.
type Response struct {
	// http status test. Example "200 OK"
	Status string

	// http status code. Example 200
	StatusCode int

	// interface that groups the basic Read and Close methods
	body io.ReadCloser

	// the request that was received by a server or to be sent by a client
	Request *http.Request
}

// GetBody
// 	Returns the response body in its i/o reader stream.
func (r *Response) GetBody() io.Reader {
	return r.body
}

// GetBodyBytes
// 	Reads the response body i/o and converts it to a byte array.
func (r *Response) GetBodyBytes() []byte {
	bodyBytes, _ := ioutil.ReadAll(r.body)
	return bodyBytes
}

// GetBodyString
// 	Parses the response body i/o interface into a string.
func (r *Response) GetBodyString() string {
	bodyBytes, _ := ioutil.ReadAll(r.body)
	return string(bodyBytes)
}
