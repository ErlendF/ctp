package riot

import (
	"ctp/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
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
	logrus.Debugf("GetLolPlaytime")
	logrus.Debugf("reg: %+v", reg)
	if reg == nil || reg.SummonerRegion == "" || reg.AccountID == "" {
		return nil, fmt.Errorf("missing summonerinfo")
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
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("StatusCode", resp.StatusCode).Warn()
		return nil, fmt.Errorf(string(resp.StatusCode))
	}

	var matches models.MatchList

	err = json.NewDecoder(resp.Body).Decode(&matches)

	var duration int
	duration = matches.TotalGames * 35 / 60

	game := &models.Game{Name: "LeagueOfLegends", Time: duration}
	return game, nil
}

//ValidateSummoner validates the summoner
func (r *Riot) ValidateSummoner(reg *models.SummonerRegistration) (*models.SummonerRegistration, error) {
	if reg == nil {
		return nil, fmt.Errorf("Nil registration")
	}
	var regions = []string{"RU", "KR", "BR1", "OC1", "JP1", "NA1", "EUN1", "EUW1", "TR1", "LA1", "LA2"}
	// var err error

	validRegion := false

	//Validating region
	for _, r := range regions {
		if r == reg.SummonerRegion {
			validRegion = true
		}
	}

	if !validRegion {
		logrus.WithField("SummonerRegion", reg.SummonerRegion).Warnf("invalid region")
		return nil, fmt.Errorf("invalid region")
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
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("SummonerName", reg.SummonerName).Warnf("invalid SummonerName")
		return nil, fmt.Errorf(string(resp.StatusCode))
	}

	var tmpReg models.SummonerRegistration

	err = json.NewDecoder(resp.Body).Decode(&tmpReg)

	reg.AccountID = tmpReg.AccountID

	return reg, nil
}
