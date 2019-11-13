package models

//Valve interface defines all methods which should be provided by valve
type Valve interface {
	GetValvePlaytime(ID string) ([]Game, error)
}

//ValveResp is used for testing
type ValveResp struct {
	Response ValveResponse `json:"response"`
}
type ValveGames struct {
	Appid                    int    `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
	ImgIconURL               string `json:"img_icon_url"`
	ImgLogoURL               string `json:"img_logo_url"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
	PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
	PlaytimeMacForever       int    `json:"playtime_mac_forever"`
	PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
	Playtime2Weeks           int    `json:"playtime_2weeks,omitempty"`
}
type ValveResponse struct {
	GameCount int     `json:"game_count"`
	Games     []ValveGames `json:"games"`
}