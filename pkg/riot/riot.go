package riot

import (
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

//Riot is a struct which contains everything necessary to handle a request related to riot
type Riot struct {
	models.Client
	apiKey string
}

//New returns a new riot instance
func New(client models.Client, apiKey string) *Riot {
	r := &Riot{apiKey: apiKey}
	r.Client = client

	return r
}

//GetRiotPlaytime gets playtime on League of Legends
func (r *Riot) GetRiotPlaytime(reg *models.SummonerRegistration) (*models.Game, error) {
	if reg == nil || reg.SummonerRegion == "" || reg.AccountID == "" {
		return nil, errors.New("missing summonerinfo")
	}

	URL := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v4/matchlists/by-account/%s?beginIndex=99999", reg.SummonerRegion, reg.AccountID)

	formatURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, formatURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", r.apiKey)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, models.NewAPIErr(err, "Riot")
	}

	defer resp.Body.Close()

	if err = models.CheckStatusCode(resp.StatusCode, "Riot", "invalid username or region for League of Legends"); err != nil {
		return nil, err
	}

	var matches models.MatchList

	err = json.NewDecoder(resp.Body).Decode(&matches)

	if err != nil {
		return nil, models.NewAPIErr(err, "Riot")
	}

	duration := matches.TotalGames * 35 / 60

	game := &models.Game{Name: "LeagueOfLegends", Time: duration}

	return game, nil
}

//ValidateSummoner validates the summoner
func (r *Riot) ValidateSummoner(reg *models.SummonerRegistration) (*models.SummonerRegistration, error) {
	if reg == nil {
		return nil, errors.New("nil summoner registration")
	}

	var regions = []string{"RU", "KR", "BR1", "OC1", "JP1", "NA1", "EUN1", "EUW1", "TR1", "LA1", "LA2"}

	if !models.Contains(regions, reg.SummonerRegion) {
		return nil, models.NewReqErrStr(fmt.Sprintf("invalid summoner region: %s", reg.SummonerRegion), "invalid region for League of Legends")
	}

	//Validating name
	URL := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", reg.SummonerRegion, reg.SummonerName)

	formatURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, formatURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", r.apiKey)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, models.NewAPIErr(err, "Riot")
	}
	defer resp.Body.Close()

	if err = models.AccValStatusCode(resp.StatusCode, "Riot", "invalid username for League of Legends"); err != nil {
		return nil, err
	}

	var tmpReg models.SummonerRegistration

	err = json.NewDecoder(resp.Body).Decode(&tmpReg)
	if err != nil {
		return nil, models.NewAPIErr(err, "Riot")
	}

	reg.AccountID = tmpReg.AccountID

	return reg, nil
}
