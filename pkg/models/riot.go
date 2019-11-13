package models

//Riot interface defines all methods which should be provided by riot
type Riot interface {
	GetRiotPlaytime(reg *SummonerRegistration) (*Game, error)
	ValidateSummoner(reg *SummonerRegistration) (*SummonerRegistration, error)
}

//MatchList should have a meaningfull comment - TODO
type MatchList struct {
	Matches    []Matches `json:"matches"`
	EndIndex   int       `json:"endIndex"`
	StartIndex int       `json:"startIndex"`
	TotalGames int       `json:"totalGames"`
}

//Matches should have a meaningfull comment - TODO
type Matches struct {
	Lane       string `json:"lane"`
	GameID     int64  `json:"gameId"`
	Champion   int    `json:"champion"`
	PlatformID string `json:"platformId"`
	Timestamp  int64  `json:"timestamp"`
	Queue      int    `json:"queue"`
	Role       string `json:"role"`
	Season     int    `json:"season"`
}

//SummonerRegistration contains the necessary information to register a summoner (league of legends account)
type SummonerRegistration struct {
	SummonerName   string `json:"summonerName" firestore:"summonerName"`
	SummonerRegion string `json:"summonerRegion" firestore:"summonerRegion"`
	AccountID      string `json:"accountId" firestore:"accountId"`
}
