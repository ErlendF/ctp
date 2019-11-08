package server

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler, apiVer string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.testHandler)
	r.HandleFunc("/riot", h.riotHandler)
	r.HandleFunc("/valve", h.valveHandler)

	return r
}
