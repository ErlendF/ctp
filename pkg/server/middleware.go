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

// The provided key must be comparable and should not be of type string
// or any other built-in type to avoid collisions between packages using context.
// Users of WithValue should define their own types for keys. - https://golang.org/pkg/context/#WithValue
type ctxKey string

func newMiddleware(val models.TokenValidator) *middleware {
	return &middleware{val}
}

func (m *middleware) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("auth")
		token := r.Header.Get("Authorization")
		id, err := m.ValidateToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		k := ctxKey("id")
		ctx := context.WithValue(r.Context(), k, id)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("logger!")
		logrus.WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Debugf("Request received")
		next.ServeHTTP(w, r)
	})
}