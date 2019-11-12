package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler) *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1/").Subrouter()
	g := s.Methods(http.MethodGet).Subrouter()

	g.HandleFunc("/", h.testHandler).Name("root")
	g.HandleFunc("/user/{id}", h.userHandler).Name("getUser")
	g.HandleFunc("/login", h.login).Name("login")
	g.HandleFunc("/authcallback", h.authCallbackHandler).Name("authCallback")

	// g.HandleFunc("/register/blizz", h.regshit).Name("reg")
	// g.HandleFunc("/register/valve", h.regshit).Name("reg")

	s.HandleFunc("/register/league", h.regLeague).Methods(http.MethodPost).Name("regLeague")
	s.HandleFunc("/user/{id}", h.updateUser).Methods(http.MethodPost).Name("updateUser")

	//catch all
	r.PathPrefix("/").HandlerFunc(h.notImplemented)

	return r
}
