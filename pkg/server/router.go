package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler, mw *middleware) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(h.notFound)

	get := r.PathPrefix("/api/v1/").Methods(http.MethodGet).Subrouter()

	get.HandleFunc("/", h.testHandler).Name("root")
	get.HandleFunc("/user/{id}", h.userHandler).Name("getUserByID")
	// get.HandleFunc("/user/{username:[a-zA-Z0-9 ]{1,15}}", h.userHandler).Name("getUserByID")
	get.HandleFunc("/login", h.login).Name("login")
	get.HandleFunc("/authcallback", h.authCallbackHandler).Name("authCallback")

	auth := r.PathPrefix("/api/v1/").Subrouter()
	auth.HandleFunc("/register/league", h.regLeague).Methods(http.MethodPost).Name("regLeague")
	auth.HandleFunc("/user", h.updateUser).Methods(http.MethodPost).Name("updateUser")
	auth.HandleFunc("/user", h.getUser).Methods(http.MethodGet).Name("getUser")

	get.Use(mw.log)
	auth.Use(mw.auth, mw.log)

	return r
}
