package models

//User contains all relevant information about the user
type User struct {
	ID            string               `json:"-"`
	Name          string               `json:"username"`
	TotalGameTime int                  `json:"totalPlayTime"`
	Lol           SummonerRegistration `json:"lol"`
	Valve         string               `json:"valve"`
	Games         []Game               `json:"games"`
}

//Game part of User
type Game struct {
	Name string `json:"game"`
	Time int   `json:"playTime"`
}
