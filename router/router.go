/*
Package router implements a server with a router and two endpoints used
for readiness and liveliness checks.

All functions that require a context to be passed should be given one from
the service handler request to correctly handle cancellations.
*/
package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jobaldw/shared/v2/client"
)

// The response payload in the form of an error.
type Error struct {
	Err string `json:"error"`
}

// Basic response fields.
type Message struct {
	ID  interface{} `json:"id,omitempty"`
	MSG string      `json:"msg,omitempty"`
}

// New
// creates a new http server and mux router with two endpoints for
// readiness and liveliness.
//
// By default a client's readiness and liveliness endpoints are
// "/health" and "/ready", but this can be overwritten using the paths
// variadic parameter.
func New(port int, clients map[string]client.Client, paths ...string) (http.Server, *mux.Router) {
	r := mux.NewRouter()
	livePath, readyPath := "", ""

	switch len(paths) {
	case 1:
		livePath = "/" + paths[0]
	case 2:
		livePath, readyPath = "/"+paths[0], "/"+paths[1]
	default:
		livePath, readyPath = "/health", "/ready"
	}

	r.HandleFunc(livePath, health()).Methods(http.MethodGet)
	r.HandleFunc(readyPath, ready(clients)).Methods(http.MethodGet)

	return http.Server{Addr: fmt.Sprintf(":%d", port)}, r
}

// RespondError
// writes a client error reply based on its passed in encoding func.
// (e.g. "json.Marshal()" or "xml.Marshal()")
func RespondError(w http.ResponseWriter, encoding func(v any) ([]byte, error), code int, err error) {
	e := Error{Err: err.Error()}
	Respond(w, encoding, code, e)
}

// Respond
// writes a client message reply based on its passed in encoding func
// (e.g. "json.Marshal()" or "xml.Marshal()").
func Respond(w http.ResponseWriter, encoding func(v any) ([]byte, error), code int, payload interface{}) {
	response, _ := encoding(payload)

	// If the Header does not contain a "Content-Type" key, Write adds a "Content-Type" set
	// to the result of passing the initial 512 bytes of written data to DetectContentType().
	w.Header().Set("Content-Type", "application/json")

	// If WriteHeader is not called, Write() calls WriteHeader(http.StatusOK).
	w.WriteHeader(code)

	w.Write(response) // nolint:errcheck
}

// Vars
// returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

/********** helper functions **********/

// health
// responses with "200 OK" status whenever called.
func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		Respond(w, json.Marshal, code, Message{MSG: http.StatusText(code)})
	}
}

// ready
// calls every dependent clients ready check endpoint.
func ready(clients map[string]client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		for _, client := range clients {
			isReady, err := client.IsReady(ctx)
			if err != nil {
				RespondError(w, json.Marshal, http.StatusBadGateway, err)
				return
			}

			if !isReady {
				Respond(w, json.Marshal, http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
				return
			}
		}

		Respond(w, json.Marshal, http.StatusOK, Message{MSG: "Ready"})
	}
}
