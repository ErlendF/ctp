package models

//Valve interface defines all methods which should be provided by valve
type Valve interface {
	GetValvePlaytime(ID string) (*ValveResp, error)
}

//ValveResp is used for testing
type ValveResp struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			Appid                  int `json:"appid"`
			PlaytimeForever        int `json:"playtime_forever"`
			PlaytimeWindowsForever int `json:"playtime_windows_forever"`
			PlaytimeMacForever     int `json:"playtime_mac_forever"`
			PlaytimeLinuxForever   int `json:"playtime_linux_forever"`
			Playtime2Weeks         int `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}
