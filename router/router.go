package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Resp from client
type Resp struct {
	ID      interface{} `json:"id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`

	Status string `json:"status,omitempty"`
	MSG    string `json:"msg,omitempty"`
	ERR    string `json:"error,omitempty"`
}

// Health check
func Health(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Resp{Status: "Up", MSG: fmt.Sprintf("%s is healthy", name)}
		Response(w, http.StatusOK, resp)
	}
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

// Response writes to client
func Response(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return errors.WithMessage(err, "could not marshal payload")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		return errors.WithMessage(err, "could not write response")
	}

	return err
}
