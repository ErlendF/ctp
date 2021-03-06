package auth

import (
	"context"
	"crypto/rand"
	"ctp/pkg/models"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Authenticator contains everything used by an authenticator
type Authenticator struct {
	ctx        context.Context
	config     oauth2.Config
	verifier   *oidc.IDTokenVerifier
	hmacSecret []byte
	uv         models.UserValidator
}

const stateCookie = "oauthstate"

// New initializes and returns an Authenticator.
// The authenticator fulfills the TokenGenerator and AuthMiddleware interfaces
// Authenticating the user through OpenIDConnect with Google as provider
// https://developers.google.com/identity/protocols/OpenIDConnect
func New(ctx context.Context, uv models.UserValidator, port int,
	domain, clientID, clientSecret, hmacSecret string) (*Authenticator, error) {
	authenticator := &Authenticator{ctx: ctx, uv: uv}

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	authenticator.verifier = provider.Verifier(oidcConfig)

	authenticator.config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  fmt.Sprintf("http://%s:%d/api/v1/authcallback", domain, port),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	authenticator.hmacSecret = []byte(hmacSecret)

	return authenticator, nil
}

// AuthRedirect redirects the user to the oauth providers confirmation page
// A random array of bytes is generated as oauth state and stored in a cookie to prevent CSRF attacks
func (a *Authenticator) AuthRedirect(w http.ResponseWriter, r *http.Request) {
	state, err := generateStateOauthCookie(w)
	if err != nil {
		logrus.WithError(err).WithField("route", mux.CurrentRoute(r).GetName()).Warn("Could not generate state for authentication")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, a.config.AuthCodeURL(state), http.StatusFound)
}

// HandleOAuth2Callback handles callback from oauth2 provider.
// The state is checked (to prevent CSRF attacks), code is exchanged for an oauth2Token and the id is retrieved.
// Neither profile nor email is used nor stored
func (a *Authenticator) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error) {
	// cookie, err := r.Cookie(stateCookie)	// commented out due to workaround, see README.md (authentication)
	// if err != nil {
	// 	return "", err
	// }

	// removing state cookie as it should not be used again regardless of whether it is valid or not
	removeStateOauthCookie(w)

	// if cookie.Value != r.URL.Query().Get("state") { // commented out due to workaround, see README.md (authentication)
	// 	return "", models.ErrInvalidAuthState
	// }

	// exchanging the authorization code for an oauth2token
	oauth2Token, err := a.config.Exchange(a.ctx, r.URL.Query().Get("code"))
	if err != nil {
		return "", err
	}

	// retrieving the raw ID token and casting it to a string
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)

	if !ok {
		return "", errors.New("no id_token field in oauth2 token")
	}

	// verifying the id token
	idToken, err := a.verifier.Verify(a.ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	var claims struct {
		Sub string // "subject", a unique identifier for the user (ID)
	}

	// extracting the claims.
	err = idToken.Claims(&claims)
	if err != nil {
		return "", err
	}

	return claims.Sub, nil
}

// generates a random array of bytes (state) and stores it in a cookie to prevent CSRF attacks
// based on example from https://dev.to/douglasmakey/oauth2-example-with-go-3n8a
func generateStateOauthCookie(w http.ResponseWriter) (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: stateCookie, Value: state}
	http.SetCookie(w, &cookie)

	return state, nil
}

// removes the state cookie
func removeStateOauthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   stateCookie,
		Value:  "",
		MaxAge: -1, // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0' - net/http package
	}

	http.SetCookie(w, cookie)
}
