package models

// Jagex interface defines all methods which should be provided by jagex
type Jagex interface {
	GetRSPlaytime(rsAcc *RunescapeAccount) (*Game, error)
	ValidateRSAccount(rsAcc *RunescapeAccount) error
}

type RunescapeAccount struct {
	Username    string `json:"username" firebase:"username"`
	AccountType string `json:"accountType" firebase:"accountType"`
	TotalLevel  int    `json:"totalLevel" firebase:"totalLevel"`
	TotalXP     int    `json:"totalXP" firebase:"totalXP"`
}
