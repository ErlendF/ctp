package models

// Database contains all functions a database should provide
type Database interface {
	SetUser(info *UserInfo) error
	GetUser(id string) (*UserInfo, error)
}
