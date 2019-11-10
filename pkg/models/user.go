package models

//UserInfo contains all relevant information about the user
type UserInfo struct {
	Name          string `json:"username"`
	TotalGameTime int    `json:"totalPlayTime"`
}
