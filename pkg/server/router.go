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
	g.HandleFunc("/riot", h.riotHandler).Name("riot")
	g.HandleFunc("/valve/{id:[0-9]+}", h.valveHandler).Name("valve")
	g.HandleFunc("/blizzard", h.blizzardHandler).Name("blizzard")
	g.HandleFunc("/user/{id}", h.userHandler).Name("Userinfo")
	g.HandleFunc("/login", h.login).Name("login")
	g.HandleFunc("/authcallback", h.authCallback).Name("authCallback")

	s.HandleFunc("/user/{id}", h.updateUser).Methods(http.MethodPost).Name("updateUser")

	return r
}
