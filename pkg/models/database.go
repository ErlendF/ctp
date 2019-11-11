package models

// Database contains all functions a database should provide
type Database interface {
	SetUser(info *User) error
	GetUser(id string) (*User, error)
}
