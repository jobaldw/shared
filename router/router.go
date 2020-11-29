package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jobaldw/shared/client"

	"github.com/gorilla/mux"
)

var appName string

// New mux router
func New(app string, oneClient client.Client, clients map[string]client.Client) *mux.Router {
	appName = app

	r := mux.NewRouter()
	r.HandleFunc("/health", health()).Methods(http.MethodGet)
	if clients == nil {
		clients := make(map[string]client.Client)
		clients["client"] = oneClient
	}

	r.HandleFunc("/ready", ready(clients)).Methods(http.MethodGet)

	return r
}

// RespondWithJSON to client
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
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

// RespondWithError to client
func RespondWithError(w http.ResponseWriter, code int, e error) error {
	if err := RespondWithJSON(w, code, struct {
		Err string `json:"error"`
	}{e.Error()}); err != nil {
		return err
	}
	return nil
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
func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, struct {
			Status string `json:"status"`
			MSG    string `json:"message"`
		}{"up", fmt.Sprintf("%s is healthy", appName)})
	}
}

func ready(clients map[string]client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, client := range clients {
			resp, err := client.Get(client.Health, nil)
			if err != nil {
				RespondWithError(w, http.StatusNotFound, err)
				return
			}

			if !IsSuccessful(resp.Code) {
				RespondWithError(w, http.StatusServiceUnavailable, fmt.Errorf("%d, %s", resp.Code, resp.Body.String))
				return
			}
		}

		RespondWithJSON(w, http.StatusOK, struct {
			Status string `json:"status"`
			MSG    string `json:"message"`
		}{"up", fmt.Sprintf("%s is ready", appName)})
	}
}
