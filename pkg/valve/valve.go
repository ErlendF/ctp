package valve

import (
	"ctp/pkg/models"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

//Valve is a struct which contains everything necessary to handle a request related to valve
type Valve struct {
	models.Getter
	apiKey string
}

//New returns a new valve instance
func New(getter models.Getter, apiKey string) *Valve {
	v := &Valve{apiKey: apiKey}
	v.Getter = getter

	return v
}

//GetValvePlaytime gets playtime on steam for specified game
func (v *Valve) GetValvePlaytime(ID string) ([]models.Game, error) {
	resp, err := v.Get(fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&format=json&steamid=%s&include_appinfo=true", v.apiKey, ID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var valvegames models.ValveResp
	err = json.NewDecoder(resp.Body).Decode(&valvegames)

	if err != nil {
		return nil, err
	}

	logrus.Debugf("GetSteamPlaytime")

	var games []models.Game

	for _ , game := range valvegames.Response.Games {
		var tmpGame models.Game
		tmpGame.Name = game.Name
		tmpGame.Time = game.PlaytimeForever

		games = append(games, tmpGame)
	}

	return games, nil
}
