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

// Handler embedds the models.UserManager interface
// which contains all functions to manage a user
type handler struct {
	models.UserManager
}

//newHandler returns handler
func newHandler(um models.UserManager) *handler {
	return &handler{um}
}

func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("testHandler!")
	//debug start
	h.JohanTestFunc()
	fmt.Fprintf(w, "Test handler!")
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

func (h *handler) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := h.AuthCallback(w, r)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Error getting token")

		//returning errorcode based on error
		switch {
		case err.Error() == models.InvalidAuthState:
			http.Error(w, fmt.Sprintf("%s: invalid state", http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}
	respondPlain(w, r, resp)
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

func (h *handler) regLeague(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("regLeague")

	var regInfo models.SummonerRegistration

	err := json.NewDecoder(r.Body).Decode(&regInfo)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	err = h.RegisterLeague(&regInfo)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	//write ok?
}

func (h *handler) notImplemented(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("notImplemented!")
	fmt.Fprintf(w, "Not implemented!")
}
