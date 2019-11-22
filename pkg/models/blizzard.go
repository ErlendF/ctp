package models

// Blizzard interface defines all methods which should be provided by blizzard
type Blizzard interface {
	GetBlizzardPlaytime(*Overwatch) (*Game, error)
	ValidateBattleUser(overwatch *Overwatch) error
}

// Overwatch struct contains users battle tag and total playtime
type Overwatch struct {
	BattleTag string `json:"battleTag" firebase:"battleTag"`
	Platform  string `json:"platform" firebase:"platform"`
	Region    string `json:"region" firebase:"region"`
}
