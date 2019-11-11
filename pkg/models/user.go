package models

//User contains all relevant information about the user
type User struct {
	ID            string `json:"-"`
	Token         string `json:"token"`
	Name          string `json:"username"`
	TotalGameTime int    `json:"totalPlayTime"`
}
