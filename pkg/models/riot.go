package models

import "time"

//Riot interface defines all methods which should be provided by riot
type Riot interface {
	GetLolPlaytime() (*time.Duration, error)
}
