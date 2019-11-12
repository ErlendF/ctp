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
	//debug start

	tmpGame := models.Game{
		Name: "League",
		Time: 12,
	}

	tmpUser := models.User{
		ID:            "117575669351657432712",
		Token:         "",
		Name:          "Johan",
		TotalGameTime: 12,
		Games:         nil,
	}

	tmpUser.Games = append(tmpUser.Games, tmpGame)
	//debug end

	err := h.SetUser(&tmpUser)
	if err != nil {
		logrus.WithError(err).Debugf("Test failed!")
	}

	game, _ := h.GetRiotPlaytime()
	err = h.UpdateGame(tmpUser.ID, game)
	if err != nil {
		logrus.WithError(err).Warnf("Update game failed!")
	}

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

	resp, err := h.GetBlizzardPlaytime("test", "test", "test")
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) userHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Debugf("Request received")

	resp, err := h.GetUser(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	h.Redirect(w, r)
}

func (h *handler) authCallback(w http.ResponseWriter, r *http.Request) {
	id, err := h.HandleOAuth2Callback(w, r)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Error getting status")

		//returning errorcode based on error
		switch {
		case err.Error() == models.InvalidAuthState:
			http.Error(w, fmt.Sprintf("%s: invalid state", http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	err = h.SetUser(&models.User{ID: id})
	if err != nil {
		logRespond(w, r, err)
		return
	}

	token, err := h.GetNewToken(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	test, err := h.ValidateToken(token)
	if err != nil {
		logrus.Debugf("Failed!")
		logRespond(w, r, err)
		return
	}
	if test != id {
		logrus.Debugf("Failed, not equal!")
		logRespond(w, r, err)
		return
	}

	respondPlain(w, r, token)
}

func (h *handler) updateUser(w http.ResponseWriter, r *http.Request) {

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

func respondPlain(w http.ResponseWriter, r *http.Request, resp string) {
	_, err := fmt.Fprintf(w, resp)
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

func (h *handler) notImplemented(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("notImplemented!")
	fmt.Fprintf(w, "Not implemented!")
}
