package models

// Database contains all functions a database should provide
type Database interface {
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByName(name string) (*User, error)
	UpdateUser(user *User) error
	UpdateGames(user *User) error
	SetUsername(user *User) error
	OverwriteUser(info *User) error
}
