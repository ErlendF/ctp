package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler, mw *middleware) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(h.notImplemented)

	get := r.PathPrefix("/api/v1/").Methods(http.MethodGet).Subrouter()

	get.HandleFunc("/", h.testHandler).Name("root")
	get.HandleFunc("/user/{id}", h.userHandler).Name("getUser")
	get.HandleFunc("/login", h.login).Name("login")
	get.HandleFunc("/authcallback", h.authCallbackHandler).Name("authCallback")

	// get.HandleFunc("/register/blizz", h.regshit).Name("reg")
	// get.HandleFunc("/register/valve", h.regshit).Name("reg")

	auth := r.PathPrefix("/api/v1/").Subrouter()
	auth.HandleFunc("/register/league", h.regLeague).Methods(http.MethodPost).Name("regLeague")
	auth.HandleFunc("/user/{id}", h.updateUser).Methods(http.MethodPost).Name("updateUser")

	get.Use(mw.log)
	auth.Use(mw.auth, mw.log)

	return r
}
