package models

// Blizzard interface defines all methods which should be provided by blizzard
type Blizzard interface {
	GetBlizzardPlaytime(*Overwatch) (*Game, error)
	ValidateBattleUser(overwatch *Overwatch) error
}

// BlizzardResp struct retrieves only time played in Overwatch
type BlizzardResp struct {
	CompetitiveStats struct {
		CareerStats CareerStats `json:"careerStats"`
	} `json:"competitiveStats"`
	QuickPlayStats struct {
		CareerStats CareerStats `json:"careerStats"`
	} `json:"quickPlayStats"`
}

// CareerStats includes time played for all heroes by game
type CareerStats struct {
	AllHeroes struct {
		Game struct {
			TimePlayed string `json:"timePlayed"`
		} `json:"game"`
	} `json:"allHeroes"`
}

// Overwatch struct contains users battle tag and total playtime
type Overwatch struct {
	BattleTag string `json:"battleTag" firebase:"battleTag"`
	Platform  string `json:"platform" firebase:"platform"`
	Region    string `json:"region" firebase:"region"`
}
