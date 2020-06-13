package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jobaldw/shared/client"

	"github.com/gorilla/mux"
)

// const status'
var (
	statusDown = "down"
	statusUp   = "up"
)

// Resp struct
type Resp struct {
	ID      interface{} `json:"id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`

	Status string `json:"status,omitempty"`
	MSG    string `json:"msg,omitempty"`
	ERR    string `json:"error,omitempty"`
}

// New mux router
func New(app string, clients map[string]client.Client) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", health(app)).Methods(http.MethodGet)
	r.HandleFunc("/ready", ready(app, clients)).Methods(http.MethodGet)

	return r
}

// Response to client
func Response(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%s, %s", err, "could not marshal payload")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		return fmt.Errorf("%s, %s", err, "could not write response")
	}

	return err
}

// IsSuccessful response code
func IsSuccessful(code int) bool {
	if (code > 199) && (code < 300) {
		return true
	}
	return false
}

// Vars from request
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// Helper functions
func health(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Resp{Status: statusUp, MSG: fmt.Sprintf("%s is healthy", name)}
		Response(w, http.StatusOK, resp)
	}
}

func ready(name string, clients map[string]client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Resp{Status: statusDown, MSG: fmt.Sprintf("%s is not ready", name)}

		for clientKey, client := range clients {
			res, err := client.Get(client.Health)
			if err != nil {
				resp.ERR = fmt.Sprintf("could not check health of %s, %s", clientKey, err)
				Response(w, http.StatusNotFound, resp)
				return
			}

			if !IsSuccessful(res.Code) {
				Response(w, http.StatusServiceUnavailable, resp)
				return
			}
		}

		resp.Status = statusUp
		resp.MSG = fmt.Sprintf("%s is ready", name)
		Response(w, http.StatusOK, resp)
	}
}
