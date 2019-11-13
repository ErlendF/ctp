package models

// Database contains all functions a database should provide
type Database interface {
	CreateUser(user *User) error
	SetUser(info *User) error
	SetUsername(user *User) error
	GetUser(id string) (*User, error)
	GetUserByName(name string) (*User, error)
	UpdateGame(userID string, game *Game) error
	UpdateUser(user *User) error
}
