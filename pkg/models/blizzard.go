package models

import "time"

// Blizzard interface defines all methods which should be provided by blizzard
type Blizzard interface {
	GetBlizzardPlaytime(*Overwatch) (*Overwatch, error)
}

// BlizzardResp struct retrieves only time played in Overwatch
type BlizzardResp struct {
	CompetitiveStats	struct{
		CareerStats    struct {
			AllHeroes struct {
				Game struct {
					TimePlayed string `json:"timePlayed"`
				} `json:"game"`
			} `json:"allHeroes"`
		} `json:"careerStats"`
	} `json:"competitiveStats"`
	QuickPlayStats		  struct{
		CareerStats      struct {
			AllHeroes	struct{
				Game   struct{
					TimePlayed	string	`json:"timePlayed"`
				} `json:"game"`
			} `json:"allHeroes"`
		} `json:"careerStats"`
	} `json:"quickPlayStats"`
}

// Overwatch struct contains users BATTLE-ID and total playtime
type Overwatch struct {
	ID         string          `json:"id" firebase:"id"`
	Platform   string          `json:"platform" firebase:"platform"`
	Region     string          `json:"region" firebase:"region"`
	Playtime   time.Duration   `json:"playtime" firebase:"playtime"`
}
