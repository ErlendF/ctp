package models

import "time"

//Riot interface defines all methods which should be provided by riot
type Riot interface {
	GetRiotPlaytime() (*time.Duration, error)
}

//MatchList should have a meaningfull comment - FIX
type MatchList struct {
	Matches    []Matches `json:"matches"`
	EndIndex   int       `json:"endIndex"`
	StartIndex int       `json:"startIndex"`
	TotalGames int       `json:"totalGames"`
}

//Matches should have a meaningfull comment - FIX
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
