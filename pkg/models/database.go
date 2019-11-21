package models

// Database contains all functions a database should provide
type Database interface {
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByName(name string) (*User, error)
	UpdateUser(user *User) error
	UpdateGames(user *User) error
	SetUsername(user *User) error
	DeleteUser(id string) error
}

// UserValidator defines the function "IsUser", which checks
// whether or not the given id is a valid user stored in the database
type UserValidator interface {
	IsUser(id string) (bool, error)
}
