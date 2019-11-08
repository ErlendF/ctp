package models

import "time"

//Valve interface defines all methods which should be provided by valve
type Valve interface {
	GetValvePlaytime(game string) (*time.Duration, error)
}
