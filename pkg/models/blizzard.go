package models

import "time"

// Blizzard interface defines all methods which should be provided by blizzard
type Blizzard interface {
	GetBlizzardPlaytime(*Overwatch) (*Overwatch, error)
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

//CareerStats includes timeplayed for all heroes by game
type CareerStats struct {
	AllHeroes struct {
		Game struct {
			TimePlayed string `json:"timePlayed"`
		} `json:"game"`
	} `json:"allHeroes"`
}

// Overwatch struct contains users BATTLE-ID and total playtime
type Overwatch struct {
	ID       string        `json:"id" firebase:"id"`
	Platform string        `json:"platform" firebase:"platform"`
	Region   string        `json:"region" firebase:"region"`
	Playtime time.Duration `json:"playtime" firebase:"playtime"`
}
