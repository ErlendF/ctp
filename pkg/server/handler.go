package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"ctp/pkg/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler embedds the models.UserManager interface
// which contains all functions to manage a user
type handler struct {
	models.UserManager
}

// newHandler returns a new handler
func newHandler(um models.UserManager) *handler {
	return &handler{um}
}

func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("testHandler!")

	h.JohanTestFunc()
	fmt.Fprintf(w, "Test handler!")
}

// check if user is public
func (h *handler) getPublicUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	resp, err := h.GetUserByName(username)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	resp.Public = false // as the user has to be public, this information is not useful

	respond(w, r, resp)
}

// log in redirect
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	h.Redirect(w, r)
}

func (h *handler) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := h.AuthCallback(w, r)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("error getting token")

		// returning errorcode based on error
		switch {
		case errors.Is(err, models.ErrInvalidAuthState):
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		err = models.NewReqErr(err, "invalid request body")
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

	resp, err := h.GetUserByID(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respond(w, r, resp)
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	err = h.DeleteUser(id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respondPlain(w, r, "Success")
}

func (h *handler) updateKey(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = models.NewReqErr(err, "invalid request body")
		logRespond(w, r, err)
		return
	}

	err = h.UpdateRiotAPIKey(string(body), id)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	respondPlain(w, r, "Success")
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
	_, err := fmt.Fprint(w, resp)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Could not encode response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func logRespond(w http.ResponseWriter, r *http.Request, err error) {
	logrus.WithField("route", mux.CurrentRoute(r).GetName()).Warn(err)

	var reqErr *models.RequestError
	var apiErr *models.ExternalAPIError
	netErr, netErrOK := err.(net.Error)

	switch {
	case errors.Is(err, models.ErrInvalidID):
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	case errors.Is(err, models.ErrNotFound):
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case errors.As(err, &reqErr):
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusBadRequest), reqErr.Response), http.StatusBadRequest)
	case errors.As(err, &apiErr):
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusBadGateway), apiErr.Respond()), http.StatusBadGateway)
	case netErrOK:
		if netErr.Timeout() {
			http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
			return
		}

		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *handler) notFound(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("request", r.RequestURI).Debug("Not found handler")
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func getID(r *http.Request) (string, error) {
	id := r.Context().Value(models.CtxKey("id"))
	idStr, ok := id.(string)

	if !ok {
		return "", models.ErrInvalidID
	}

	return idStr, nil
}
