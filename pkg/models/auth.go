package models

import "net/http"

//TokenValidator validates a token
type TokenValidator interface {
	ValidateToken(tokenString string) (string, error)
	IsUser(id string) (bool, error)
}

//TokenGenerator generates a new token
type TokenGenerator interface {
	GetNewToken(id string) (string, error)
	AuthRedirect(w http.ResponseWriter, r *http.Request)
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error)
}
