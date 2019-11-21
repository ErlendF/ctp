package models

// Jagex interface defines all methods which should be provided by jagex
type Jagex interface {
	GetRSPlaytime(username string) (*Game, error)
	ValidateRSAccount(username string) error
}
