package models

import "time"

// Blizzard interface defines all methods which should be provided by blizzard
type Blizzard interface {
	GetBlizzardPlaytime(platform, region, ID string) (*time.Duration, error)
}

// BlizzardResp struct retrieves only time played in Overwatch
type BlizzardResp struct {
	CompetitiveStats	struct{
		AllHeroes	struct{
			Game struct{
				TimePlayed	string	`json:"timePlayed"`
			} `json:"game"`
		} `json:"allHeroes"`
	} `json:"competitiveStats"`
	QuickPlayStats		struct{
		AllHeroes	struct{
			Game struct{
				TimePlayed	string	`json:"timePlayed"`
			} `json:"game"`
		} `json:"allHeroes"`
	} `json:"quickPlayStats"`
}