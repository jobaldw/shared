package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jobaldw/shared/client"
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
// 	Creates a new http server and mux router with to readiness endpoints: "/health" and "/ready".
// 	* @param port: TCP address the server will listen on
// 	* @param clients: map of clients to test readiness
// 	* @param paths: rename live and ready paths
func New(port int, clients map[string]client.Client, paths ...string) (http.Server, *mux.Router) {
	r := mux.NewRouter()
	var livePath, readyPath = "", ""

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
//	Writes a client error reply based on its passed in encoding func.
// 	* @param w: the handlers way to construct a response
// 	* @param encoding: a func to encode the payload, e.g. "json.Marshal()" or "xml.Marshal()"
// 	* @param code: HTTP response status code
// 	* @param err: response error
func RespondError(w http.ResponseWriter, encoding func(v any) ([]byte, error), code int, err error) {
	e := Error{Err: err.Error()}
	Respond(w, encoding, code, e)
}

// Respond
//	Writes a client message reply based on its passed in encoding func.
// 	* @param w: the handlers way to construct a response
// 	* @param encoding: a func to encode the payload; normally json.Marshal() or xml.Marshal()
// 	* @param code: HTTP response status code
// 	* @param payload: response body
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
// 	returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

/********** Helper functions **********/

// health
// 	Responses with "200 OK" status whenever called.
func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		Respond(w, json.Marshal, code, Message{MSG: http.StatusText(code)})
	}
}

// ready
// 	Calls every dependent clients ready check endpoint.
// 	* @param clients: map of clients to test readiness
func ready(clients map[string]client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, client := range clients {
			isReady, err := client.IsReady()
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
