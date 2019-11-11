package models

//UserManager contains all functions a usermanager is expected to provide
type UserManager interface {
	GetUser(username string) (*User, error)
	SetUser(user *User) error
}
