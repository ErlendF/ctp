package server

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler) *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1/").Subrouter()

	s.HandleFunc("/", h.testHandler).Name("root")
	s.HandleFunc("/riot", h.riotHandler).Name("riot")
	s.HandleFunc("/valve/{id:[0-9]+}", h.valveHandler).Name("valve")
	s.HandleFunc("/blizzard", h.blizzardHandler).Name("blizzard")
	s.HandleFunc("/user/{id}", h.userHandler).Name("Userinfo")
	s.HandleFunc("/login", h.login).Name("login")
	s.HandleFunc("/loginRedirected", h.redirected).Name("redirected")

	return r
}
