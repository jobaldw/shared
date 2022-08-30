package errors

import (
	"fmt"

	"github.com/jobaldw/shared/v2/client"
)

// struct for client errors
type ClientErr struct {
	Client     string `json:"client,omitempty"`
	URI        string `json:"uri,omitempty"`
	Status     string `json:"status,omitempty"`
	StatusCode int    `json:"code,omitempty"`
	Msg        string `json:"message,omitempty"`
	Err        error  `json:"error,omitempty"`
}

// NewClientErr
//
//	Creates and fills a new MongoErr struct.
//	* @param op: the operation
//	* @param col: the mongo collection
//	* @param err: an error
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
//
// Error implements the error interface.
func (ce ClientErr) Error() string {
	return fmt.Sprintf("[%s] error: %v", ce.Client, ce.Err)
}
