package server

import (
	"context"
	"ctp/pkg/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type middleware struct {
	models.TokenValidator
}

// "The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys." - https://golang.org/pkg/context/#WithValue
type ctxKey string

func newMiddleware(val models.TokenValidator) *middleware {
	return &middleware{val}
}

// auth validates received token and passes the id to handlers by request context
func (m *middleware) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		id, err := m.ValidateToken(token)
		if err != nil {
			logrus.WithError(err).Warn("Invalid authorization")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Checking whether or not the user exists in the database.
		// A user can have a valid token, but not exist in the database if they have deleted their account
		validUser, err := m.IsUser(id)
		if err != nil {
			logrus.WithError(err).Warn("error getting user from database")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !validUser {
			logrus.Warn("non-existing user with valid token tried to login")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		k := ctxKey("id")
		ctx := context.WithValue(r.Context(), k, id)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// log logs requests received with the routeName
func (m *middleware) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Debugf("Request received")
		next.ServeHTTP(w, r)
	})
}
