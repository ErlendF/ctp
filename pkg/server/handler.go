package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"ctp/pkg/models"

	"github.com/gorilla/mux"
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
	id := mux.Vars(r)["id"]

	resp, err := h.GetValvePlaytime(id)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Error getting status")

		//returning errorcode based on error
		switch {
		case strings.Contains(err.Error(), models.NonOK):
			if models.CheckNotFound(err) {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Bad gateway", http.StatusBadGateway)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	setHeader(w)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Could not encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) riotHandler(w http.ResponseWriter, r *http.Request) {
	h.GetRiotPlaytime()
	fmt.Fprintf(w, "Riot!")
}

func (h *handler) blizzardHandler(w http.ResponseWriter, r *http.Request) {
	h.GetBlizzardPlaytime("test")
	fmt.Fprintf(w, "Blizzard!")
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
