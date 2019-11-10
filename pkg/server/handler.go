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

	logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName(), "ID": id}).Debugf("Request received")

	resp, err := h.GetValvePlaytime(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) riotHandler(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Debugf("Request received")

	resp, err := h.GetRiotPlaytime()
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) blizzardHandler(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Debugf("Request received")

	resp, err := h.GetBlizzardPlaytime("test")
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) userHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Debugf("Request received")

	resp, err := h.GetUserInfo(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	h.Redirect(w, r)
}

func (h *handler) redirected(w http.ResponseWriter, r *http.Request) {
	id, err := h.HandleOAuth2Callback(w, r)
	if err != nil {
		logRespond(w, r, err)
	}

	logrus.Debugf("Sucess! %s", id)
}

func respond(w http.ResponseWriter, r *http.Request, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(resp)

	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Could not encode response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func logRespond(w http.ResponseWriter, r *http.Request, err error) {
	logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Error getting status")

	//returning errorcode based on error
	switch {
	case strings.Contains(err.Error(), models.NonOK):
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
