package models

import "time"

//Blizzard interface defines all methods which should be provided by blizzard
type Blizzard interface {
	GetBlizzardPlaytime(game string) (*time.Duration, error)
}
