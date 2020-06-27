package middleware

import (
	"fmt"
	"net/http"

	"github.com/jobaldw/shared/config"
	"github.com/jobaldw/shared/log"
	"github.com/jobaldw/shared/router"

	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/handlers"

	"gopkg.in/square/go-jose.v2"
)

var domain, identifier string

// New values
func New(authentication config.Auth0) {
	domain = authentication.Domain
	identifier = authentication.Identifier
}

// Auth0 authentication
func Auth0(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, POST, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		domain := "https://" + domain + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{identifier}, domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(r)
		if err != nil {
			msg := fmt.Sprintf("unauthorized, %s", err)
			log.Details().Errorf(msg)
			resp := router.Resp{ERR: msg}
			router.Response(w, http.StatusUnauthorized, resp)
			return
		}

		next(w, r)
	})
}

// Handler middleware options
func Handler(r http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"DELETE", "GET", "POST", "PUT"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(r)
}
