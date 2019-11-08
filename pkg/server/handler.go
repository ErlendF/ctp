package server

import (
	"fmt"
	"net/http"

	"ctp/pkg/models"

	"github.com/sirupsen/logrus"
)

// Handler embedds the models.Organizer interface
// which contains interfaces to respond to each of the routes.
// Organizer is used to simplify the passing of all interfaces to the handler.
type handler struct {
	models.Organizer
}

//newHandler returns handler
func newHandler(organizer models.Organizer) *handler {
	return &handler{organizer}
}

func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("testHandler!")
	fmt.Fprintf(w, "Success!")
}

func (h *handler) valveHandler(w http.ResponseWriter, r *http.Request) {
	h.GetValvePlaytime("test")
	fmt.Fprintf(w, "Valve!")
}

func (h *handler) riotHandler(w http.ResponseWriter, r *http.Request) {
	h.GetRiotPlaytime()
	fmt.Fprintf(w, "Riot!")
}

func (h *handler) blizzardHandler(w http.ResponseWriter, r *http.Request) {
	h.GetBlizzardPlaytime("test")
	fmt.Fprintf(w, "Blizzard!")
}
