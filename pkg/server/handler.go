package server

import (
	"fmt"
	"net/http"

	"ctp/pkg/models"
)

// Handler embedds the models.Organizer interface
// which contains interfaces to respond to each of the routes.
// Organizer is used to simplify the passing of all interfaces to the handler.
type handler struct {
	models.Organizer
}

//newHandler returns handler
func newHandler(organizer models.Organizer) *handler {
	h := &handler{organizer}
	return h
}

func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Success!")
}
