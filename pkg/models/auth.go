package models

import "net/http"

// TokenGenerator generates a new token and handles OAuth redirect and callbacks
type TokenGenerator interface {
	GetNewToken(id string) (string, error)
	AuthRedirect(w http.ResponseWriter, r *http.Request)
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error)
}

// AuthMiddleware defines the functions which an AuthMiddleware should provide
type AuthMiddleware interface {
	Auth(next http.Handler) http.Handler
}

// CtxKey is used to set the ID of a user as a value in the request context.
// "The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys." - https://golang.org/pkg/context/#WithValue
type CtxKey string
