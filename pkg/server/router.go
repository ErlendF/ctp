package server

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler, apiVer string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.testHandler).Name("root")
	r.HandleFunc("/riot", h.riotHandler).Name("riot")
	r.HandleFunc("/valve/{id:[0-9]+}", h.valveHandler).Name("valve")
	r.HandleFunc("/blizzard", h.blizzardHandler).Name("blizzard")

	return r
}
