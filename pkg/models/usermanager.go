package models

import "net/http"

// UserManager contains all functions a usermanager is expected to provide for "managing" a user
type UserManager interface {
	GetUserByID(id string) (*User, error)
	GetUserByName(username string) (*User, error)
	SetUser(user *User) error
	DeleteUser(id string, fields []string) error
	UpdateRiotAPIKey(key string) error
	UpdateGames(id string) error
	Redirect(w http.ResponseWriter, r *http.Request)
	AuthCallback(w http.ResponseWriter, r *http.Request) (string, error)
}
