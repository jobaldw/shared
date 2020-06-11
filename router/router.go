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
	ERR    string `json:"error,omitempty"`
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

// Ready state
// FIXME: add controller package to shared library
// func Ready(name string, ctrl controller.Controller) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		resp := Resp{Status: "Up"}

// 		if err := ctrl.Check(); err != nil {
// 			resp.Status = "Down"
// 			resp.ERR = err.Error()

// 			log.Entry.WithField("method", "Ready").Error(err)
// 			Response(w, http.StatusGatewayTimeout, resp)
// 			return
// 		}

// 		resp.MSG = fmt.Sprintf("%s is ready", name)

// 		log.Entry.WithFields(logrus.Fields{
// 			"method": "Ready",
// 			"status": resp.Status,
// 		},
// 		).Info(resp.MSG)
// 		Response(w, http.StatusOK, resp)
// 	}
// }
