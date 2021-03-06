package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jobaldw/shared/client"
	"github.com/jobaldw/shared/log"
	"github.com/jobaldw/shared/router"

	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/handlers"

	"gopkg.in/square/go-jose.v2"
)

var domain, identifier, clientID, clientSecret string

// New values
func New(name, id, secret string) {
	domain = "https://" + os.Getenv("A0_DOMAIN")
	identifier = "https://" + name
	clientID = os.Getenv(id)
	clientSecret = os.Getenv(secret)
}

// Auth0 authentication
func Auth0(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func(e error) {
			if e != nil {
				log.File().WithError(e)
			}
		}(err)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, POST, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + "/.well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{identifier}, domain+"/", jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := getToken()
		if err != nil {
			router.RespondWithError(w, http.StatusGatewayTimeout, err)
			return
		}
		r.Header.Add("Authorization", "Bearer "+token)

		_, err = validator.ValidateRequest(r)
		if err != nil {
			router.RespondWithError(w, http.StatusUnauthorized, err)
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

// helper function
func getToken() (string, error) {
	payload := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Audience     string `json:"audience"`
		GrantType    string `json:"grant_type"`

		AccessToken string `json:"access_token,omitempty"`
	}{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     identifier,
		GrantType:    "client_credentials",
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(payload)

	auth0, err := client.New(domain, "", 10, map[string]string{"content-type": "application/json"})
	if err != nil {
		return "", fmt.Errorf("could not create client, %s", err)
	}

	resp, err := auth0.Post("oauth/token", nil, bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		return "", fmt.Errorf("could not call %s, %s", domain, err)
	}

	if err := json.Unmarshal(resp.Body.Bytes, &payload); err != nil {
		return "", fmt.Errorf("could not unmarshal response, %s", err)
	}

	return payload.AccessToken, err
}
