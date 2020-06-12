package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
func New(app string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", health(app)).Methods(http.MethodGet)

	return r
}

func health(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Resp{Status: "Up", MSG: fmt.Sprintf("%s is healthy", name)}
		Response(w, http.StatusOK, resp)
	}
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
