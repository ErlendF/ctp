package models

// Riot interface defines all methods which should be provided by riot
type Riot interface {
	GetLolPlaytime(reg *SummonerRegistration) (*Game, error)
	ValidateSummoner(reg *SummonerRegistration) error
	UpdateKey(key string) error
}

// MatchList contains a list of matches with relevant information
type MatchList struct {
	TotalGames int `json:"totalGames"`
}

// SummonerRegistration contains the necessary information to register a summoner (league of legends account)
type SummonerRegistration struct {
	SummonerName   string `json:"summonerName" firestore:"summonerName"`
	SummonerRegion string `json:"summonerRegion" firestore:"summonerRegion"`
	AccountID      string `json:"accountId" firestore:"accountId"`
}
