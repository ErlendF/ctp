package models

// Valve interface defines all methods which should be provided by valve
type Valve interface {
	ValidateValveAccount(username string) (string, error)
	ValidateValveID(id string) error
	GetValvePlaytime(ID string) ([]Game, error)
}

// ValveResp is used for testing
type ValveResp struct {
	Response ValveResponse `json:"response"`
}

// ValveGames is used for testing
type ValveGames struct {
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

// ValveResponse is used for testing
type ValveResponse struct {
	GameCount int          `json:"game_count"`
	Games     []ValveGames `json:"games"`
}

// ValveAccount contains all information about a user relevant to Valve (steam)
type ValveAccount struct {
	ID       string `json:"id,omitempty" firestore:"id"`
	Username string `json:"username,omitempty" firestore:"username"`
}
