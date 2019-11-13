package models

//User contains all relevant information about the user
type User struct {
	ID            string               `json:"-" firestore:"id"`
	Name          string               `json:"username" firestore:"name"`
	TotalGameTime int                  `json:"totalPlayTime" firestore:"totalGameTime"`
	Lol           SummonerRegistration `json:"lol" firestore:"lol"`
	Valve         string               `json:"valve" firestore:"valve"`
	Games         []Game               `json:"games" firestore:"games"`
}

//Game part of User
type Game struct {
	Name string `json:"game" firestore:"name"`
	Time int    `json:"playTime" firestore:"time"`
}
