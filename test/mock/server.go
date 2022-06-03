package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Server() *httptest.Server {
	r := mux.NewRouter()

	r.Handle("/health", health()).Methods(http.MethodGet)
	r.Handle("/save", save()).Methods(http.MethodPost)
	r.Handle("/update", update()).Methods(http.MethodPut)
	r.Handle("/delete", delete()).Methods(http.MethodDelete)

	return httptest.NewServer(cors.Default().Handler(r))
}

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)

		response, _ := json.Marshal(http.StatusText(code))
		w.Write(response)
	}
}

func save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusCreated

		var req interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			code = http.StatusBadRequest
		}

		if req != "test payload" {
			code = http.StatusUnprocessableEntity
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		response, _ := json.Marshal(http.StatusText(code))
		w.Write(response)
	}
}

func update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK

		var req interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			code = http.StatusBadRequest
		}

		if req != "test payload" {
			code = http.StatusUnprocessableEntity
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		response, _ := json.Marshal(http.StatusText(code))
		w.Write(response)
	}
}

func delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusNoContent

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		response, _ := json.Marshal(http.StatusText(code))
		w.Write(response)
	}
}
