package server

import (
	"encoding/json"
	"fmt"
	"io"
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

const invalidID = "Invalid id"

//newHandler returns handler
func newHandler(um models.UserManager) *handler {
	return &handler{um}
}

func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("testHandler!")

	h.JohanTestFunc()
	fmt.Fprintf(w, "Test handler!")
}

func (h *handler) getPublicUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	resp, err := h.GetUserByName(username)
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
	id, err := getID(r)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logRespond(w, r, err)
		return
	}
	user.ID = id

	// ignoring fields the user should not be allowed to update manually
	user.Games = nil
	user.TotalGameTime = 0

	err = h.SetUser(&user)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respondPlain(w, r, "Success")
}

func (h *handler) updateGames(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	err = h.UpdateGames(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respondPlain(w, r, "Success")
}

func (h *handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	resp, err := h.GetUser(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
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
	logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Error")

	//returning errorcode based on error
	switch {
	case err.Error() == invalidID:
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	case strings.Contains(err.Error(), models.NonOK):
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	case err == io.EOF:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *handler) notFound(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("request", r.RequestURI).Debug("Not found handler")
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func getID(r *http.Request) (string, error) {
	id := r.Context().Value(ctxKey("id"))
	idStr, ok := id.(string)
	if !ok {
		return "", fmt.Errorf(invalidID)
	}

	return idStr, nil
}
