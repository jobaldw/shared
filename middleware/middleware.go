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

var domain, identifier string

// New values
func New(id string) {
	domain = "https://" + os.Getenv("DOMAIN")
	identifier = "https://" + id
}

// Auth0 authentication
func Auth0(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := router.Resp{}

		defer r.Body.Close()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, POST, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + "/.well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{identifier}, domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := getToken()
		if err != nil {
			resp.ERR = err.Error()
			log.Details().Errorf(resp.ERR)
			router.Response(w, http.StatusGatewayTimeout, resp)
			return
		}

		fmt.Println("Bearer " + token)

		r.Header.Add("Authorization", "Bearer "+token)

		fmt.Println(r)
		_, err = validator.ValidateRequest(r)
		if err != nil {
			resp.ERR = fmt.Sprintf("unauthorized, %s", err)
			log.Details().Errorf(resp.ERR)
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

// helper function
func getToken() (string, error) {
	payload := struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Audience     string `json:"audience"`
		GrantType    string `json:"grant_type"`

		AccessToken string `json:"access_token,omitempty"`
	}{
		ClientID:     "7V1Ql8WrvTCdIOcXKA8KzViotXXTK1ka",
		ClientSecret: "ZR1MQLAMD3PuIyP4nOptkS3jGo0-UmRRK_yK9NIOJlzlLSmOedyclhRKBVn3_0r5",
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
