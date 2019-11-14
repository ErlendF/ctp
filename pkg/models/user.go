package models

//User contains all relevant information about the user
type User struct {
	ID            string                `json:"-" firestore:"id"`
	Name          string                `json:"username,omitempty" firestore:"name"`
	TotalGameTime int                   `json:"totalPlayTime" firestore:"totalGameTime"`
	Lol           *SummonerRegistration `json:"lol" firestore:"lol"`
	Valve         string                `json:"valve" firestore:"valve"`
	Overwatch     *Overwatch            `json:"overwatch" firestore:"overwatch"`
	Games         []Game                `json:"games" firestore:"games"`
}

//Game part of User
type Game struct {
	Name    string `json:"game" firestore:"name"`
	Time    int    `json:"playTime" firestore:"time"`
	ValveID int    `json:"-" firestore:"valveID"`
}
