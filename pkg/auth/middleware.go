package auth

import (
	"context"
	"ctp/pkg/models"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Auth is a middleware that validates received token and passes the id to handlers by request context.
// If the token was invalid, or some error occurred, the request is rejected and no handler is called.
func (a *Authenticator) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			logrus.Warn("no token provided")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		id, err := a.validateToken(token)
		if err != nil {
			logrus.WithError(err).Warn("invalid authorization")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Checking whether or not the user exists in the database.
		// A user can have a valid token, but not exist in the database if they have deleted their account.
		validUser, err := a.uv.IsUser(id)
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

		k := models.CtxKey("id")
		ctx := context.WithValue(r.Context(), k, id)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
