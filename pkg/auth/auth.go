package auth

import (
	"context"
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

var state = "test"

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
		RedirectURL:  "http://localhost:8080/api/v1/loginRedirected",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return authenticator, nil
}

//Redirect redirects
func (a *Authenticator) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.config.AuthCodeURL(state), http.StatusFound)
}

//HandleOAuth2Callback handles callback from oauth2
func (a *Authenticator) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (string, error) {
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
