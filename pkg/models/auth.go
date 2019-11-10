package models

import "net/http"

//Authenticator contains all functions which needs to be provided by an authenticator
type Authenticator interface {
	Redirect(w http.ResponseWriter, r *http.Request)
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error)
}
