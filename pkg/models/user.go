package models

// User contains all relevant information about the user
type User struct {
	ID            string                `json:"-" firestore:"id"`
	Name          string                `json:"name,omitempty" firestore:"name"`
	Public        bool                  `json:"public,omitempty" firestore:"public"`
	TotalGameTime int                   `json:"totalPlayTime" firestore:"totalGameTime"`
	Lol           *SummonerRegistration `json:"lol,omitempty" firestore:"lol"`
	Valve         *ValveAccount         `json:"valve,omitempty" firestore:"valve"`
	Overwatch     *Overwatch            `json:"overwatch,omitempty" firestore:"overwatch"`
	Runescape     *RunescapeAccount     `json:"runescape,omitempty" firestore:"runescape"`
	Games         []Game                `json:"games" firestore:"games"`
}

// Game contains relevant information about a game
type Game struct {
	Name string `json:"game" firestore:"name"`
	Time int    `json:"playTime" firestore:"time"`
}
