package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jobaldw/shared/client"
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
	ERR    error  `json:"error,omitempty"`
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

		for _, v := range clients {
			res, err := v.Client.Get(v.Health)
			if err != nil {
				resp.Payload = res
				resp.MSG = fmt.Sprintf("%d", res.StatusCode)
				resp.ERR = fmt.Errorf("could not check health of %s, %s", v.URL.String(), err)
				Response(w, http.StatusNotFound, resp)
				return
			}

			if !IsSuccessful(res.StatusCode) {
				Response(w, http.StatusServiceUnavailable, resp)
				return
			}
		}

		resp.Status = statusUp
		resp.MSG = fmt.Sprintf("%s is ready", name)
		Response(w, http.StatusOK, resp)
	}
}
