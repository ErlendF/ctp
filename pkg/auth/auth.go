package auth

import (
	"context"
	"crypto/rand"
	"ctp/pkg/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

//Authenticator contains everything used by an authenticator
type Authenticator struct {
	ctx      context.Context
	mu       sync.RWMutex
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
}

const stateCookie = "oauthstate"

//New returns a new authenticator
func New(ctx context.Context, clientID string, clientSecret string) (*Authenticator, error) {
	authenticator := &Authenticator{ctx: ctx}

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
		RedirectURL:  "http://localhost:8080/api/v1/authcallback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return authenticator, nil
}

//Redirect redirects
func (a *Authenticator) Redirect(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	http.Redirect(w, r, a.config.AuthCodeURL(state), http.StatusFound)
}

//HandleOAuth2Callback handles callback from oauth2
func (a *Authenticator) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(stateCookie)
	if err != nil {
		return "", err
	}

	// removing state cookie as it should not be used again regardless of whether it is valid or not
	removeStateOauthCookie(w)
	if cookie.Value != r.URL.Query().Get("state") {
		return "", fmt.Errorf(models.InvalidAuthState)
	}

	oauth2Token, err := a.config.Exchange(a.ctx, r.URL.Query().Get("code"))
	if err != nil {
		return "", err
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("No id_token field in oauth2 token")
	}
	idToken, err := a.verifier.Verify(a.ctx, rawIDToken)
	if err != nil {
		return "", err
	}

	var claims struct {
		Sub string
	}
	err = idToken.Claims(&claims)
	if err != nil {
		return "", err
	}
	return claims.Sub, nil
}

// makes a random
// based on example from https://dev.to/douglasmakey/oauth2-example-with-go-3n8a
func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: stateCookie, Value: state}
	http.SetCookie(w, &cookie)

	return state
}

func removeStateOauthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   stateCookie,
		Value:  "",
		MaxAge: -1, // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0' - net/http package
	}

	http.SetCookie(w, cookie)
}
