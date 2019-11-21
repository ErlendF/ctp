package server

import (
	"ctp/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler, amw models.AuthMiddleware) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(h.notFound)

	get := r.PathPrefix("/api/v1").Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/", h.testHandler).Name("root")
	get.HandleFunc("/login", h.login).Name("login")
	get.HandleFunc("/authcallback", h.authCallbackHandler).Name("authCallback")
	get.HandleFunc("/user/{username:[a-zA-Z0-9 ]{1,15}}", h.getPublicUser).Name("getPublicUser")

	auth := r.PathPrefix("/api/v1/").Subrouter()
	auth.HandleFunc("/user", h.getUser).Methods(http.MethodGet).Name("getUser")
	auth.HandleFunc("/user", h.updateUser).Methods(http.MethodPost).Name("updateUser")
	auth.HandleFunc("/user", h.deleteUser).Methods(http.MethodDelete).Name("deleteUser")
	auth.HandleFunc("/updategames", h.updateGames).Methods(http.MethodPost).Name("updateGames")

	// loggin every request using the log middleware
	get.Use(log)

	// users are first authenticated using the authentication middleware (checks the "Authorization" header for valid token).
	// any successful requests are logged using the log middleware.
	// The AuthMiddleware interface is implemented by the Authenticator in the auth package.
	auth.Use(amw.Auth, log)

	return r
}
