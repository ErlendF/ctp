package valve

import (
	"ctp/pkg/models"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const getOwnedGames = "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&format=json&steamid=%s&include_appinfo=true"

// Valve is a struct which contains everything necessary to handle a request related to valve
type Valve struct {
	models.Getter
	apiKey string
}

// steamResp is used for decoding the response from steam
type steamPrivateCheck struct {
	Response struct {
		Players []struct {
			ID64           string `json:"steamid"`
			VisibilityCode int    `json:"communityvisibilitystate"`
		} `json:"players"`
	} `json:"response"`
}

// steamResp is used for decoding the response from steam
type steamResp struct {
	Response struct {
		ID64 string `json:"steamid"`
		Code int    `json:"success"`
	} `json:"response"`
}

// New returns a new valve instance
func New(getter models.Getter, apiKey string) *Valve {
	v := &Valve{apiKey: apiKey}
	v.Getter = getter

	return v
}

// ValidateValveAccount validates the steam account and returns the valve 64 bit ID
func (v *Valve) ValidateValveAccount(username string) (string, error) {
	if username == "" {
		return "", models.NewReqErrStr("invalid steam account", "invalid steam account")
	}
	resp, err := v.Get(fmt.Sprintf("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s", v.apiKey, username))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = models.CheckStatusCode(resp.StatusCode, "Valve", "invalid steam username")
	if err != nil {
		return "", err
	}

	var sResp steamResp
	err = json.NewDecoder(resp.Body).Decode(&sResp)
	if err != nil {
		return "", err
	}

	if sResp.Response.Code != 1 {
		return "", models.NewReqErrStr("invalid steam account", "invalid steam account")
	}

	if !strings.HasPrefix(sResp.Response.ID64, "7656119") {
		return "", models.NewReqErrStr("invalid steam account", "invalid steam account")
	}

	err = v.checkPrivateProfile(sResp.Response.ID64)
	if err != nil {
		return "", err
	}

	return sResp.Response.ID64, nil
}

// ValidateValveID validates the 64-bit steam account id
func (v *Valve) ValidateValveID(id string) error {
	if !strings.HasPrefix(id, "7656119") {
		return models.NewReqErrStr("invalid steam id", "invalid steam id")
	}

	resp, err := v.Get(fmt.Sprintf(getOwnedGames, v.apiKey, id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = models.AccValStatusCode(resp.StatusCode, "Valve", "invalid steam username")
	if err != nil {
		return err
	}

	err = v.checkPrivateProfile(id)
	if err != nil {
		return err
	}

	return nil
}

// GetValvePlaytime gets playtime on steam for specified game
func (v *Valve) GetValvePlaytime(id string) ([]models.Game, error) {
	logrus.Debug("GetSteamPlaytime")

	resp, err := v.Get(fmt.Sprintf(getOwnedGames, v.apiKey, id))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = models.CheckStatusCode(resp.StatusCode, "Valve", "invalid steam id")
	if err != nil {
		return nil, err
	}

	var valvegames models.ValveResp
	err = json.NewDecoder(resp.Body).Decode(&valvegames)

	if err != nil {
		return nil, err
	}

	var games []models.Game

	for _, game := range valvegames.Response.Games {
		var tmpGame models.Game
		tmpGame.Name = game.Name
		tmpGame.Time = game.PlaytimeForever / 60
		tmpGame.ValveID = game.Appid

		if tmpGame.Time != 0 {
			games = append(games, tmpGame)
		}
	}

	return games, nil
}

func (v *Valve) checkPrivateProfile(id string) error {
	// Checking whether or not the passed in steam account is private
	resp, err := v.Get(fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", v.apiKey, id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var sResp steamPrivateCheck
	err = json.NewDecoder(resp.Body).Decode(&sResp)
	if err != nil {
		return err
	}
	if sResp.Response.Players[0].VisibilityCode != 3 {
		return models.NewReqErrStr("private steam account", "private steam account")
	}

	return nil
}
