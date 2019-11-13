package models

import "net/http"

//UserManager contains all functions a usermanager is expected to provide
type UserManager interface {
	GetUser(id string) (*User, error)
	GetUserByName(username string) (*User, error)
	SetUser(user *User) error
	UpdateGame(id string, game *Game) error
	Redirect(w http.ResponseWriter, r *http.Request)
	AuthCallback(w http.ResponseWriter, r *http.Request) (string, error)
	RegisterLeague(id string, reg *SummonerRegistration) error
	JohanTestFunc()
}
