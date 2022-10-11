/*
Package errors handles HTTP client and Mongo errors.
*/
package errors

import (
	"fmt"

	"github.com/jobaldw/shared/v2/client"
)

// struct for client errors
type ClientErr struct {
	// name of the client
	Client string `json:"client,omitempty"`

	// reference url that the client error originated
	URI string `json:"uri,omitempty"`

	// client error's status text
	Status string `json:"status,omitempty"`

	// client error's status code
	StatusCode int `json:"code,omitempty"`

	// client response message
	Msg string `json:"message,omitempty"`

	// client response error
	Err error `json:"error,omitempty"`
}

// NewClientErr
// creates and fills a new ClientErr struct.
func NewClientErr(client, msg string, err error, resp client.Response) ClientErr {
	return ClientErr{
		Client:     client,
		URI:        resp.Request.RequestURI,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Msg:        msg,
		Err:        err,
	}
}

// Error
// implements the error interface.
func (ce ClientErr) Error() string {
	if ce.StatusCode == 0 {
		return fmt.Sprintf("[%s] error: %v", ce.Client, ce.Err)
	}
	return fmt.Sprintf("[%s] error: %d | %v", ce.Client, ce.StatusCode, ce.Err)
}
