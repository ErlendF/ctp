package riot

import (
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Riot is a struct which contains everything necessary to handle a request related to riot
type Riot struct {
	models.Client
	apiKey string

	// strings are not thread safe. The mutex is used as a VERY simple locking mechanism to facilitate the UpdateKey hack. See readme
	mutex *sync.Mutex
}

// New returns a new riot instance
func New(client models.Client, apiKey string) *Riot {
	r := &Riot{apiKey: apiKey, mutex: &sync.Mutex{}}
	r.Client = client

	return r
}

// GetLolPlaytime gets playtime on League of Legends
func (r *Riot) GetLolPlaytime(reg *models.SummonerRegistration) (*models.Game, error) {
	if reg == nil || reg.SummonerRegion == "" || reg.AccountID == "" {
		return nil, errors.New("missing summonerinfo")
	}

	// create an URL with needed params
	URL := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v4/matchlists/by-account/%s?beginIndex=99999",
		reg.SummonerRegion, reg.AccountID)

	// ensure that the URL is correctly formatted
	formatURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	// create http request
	req, err := http.NewRequest(http.MethodGet, formatURL.String(), nil)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	// set header token to avoid getting "403 unauthorized" from API
	req.Header.Set("X-Riot-Token", r.apiKey)
	r.mutex.Unlock()

	// query riot api
	resp, err := r.Do(req)
	if err != nil {
		return nil, models.NewAPIErr(err, "Riot")
	}
	defer resp.Body.Close()

	// checks the response status code and returns appropriate error
	err = models.CheckStatusCode(resp.StatusCode, "Riot", "invalid username or region for League of Legends")
	if err != nil {
		return nil, err
	}

	// decode list of matches retrieved from API
	var matches models.MatchList
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		return nil, models.NewAPIErr(err, "Riot")
	}

	// calculate total game time, assuming average match time is 35 minutes
	duration := matches.TotalGames * 35 / 60

	// create and return expected struct
	game := &models.Game{Name: "LeagueOfLegends", Time: duration}
	return game, nil
}

// ValidateSummoner validates the summoner
func (r *Riot) ValidateSummoner(reg *models.SummonerRegistration) error {
	if reg == nil {
		return errors.New("nil summoner registration")
	}

	// Checks that the payload (reg) contains a valid region
	var regions = []string{"RU", "KR", "BR1", "OC1", "JP1", "NA1", "EUN1", "EUW1", "TR1", "LA1", "LA2"}
	if !models.Contains(regions, reg.SummonerRegion) {
		return models.NewReqErrStr(fmt.Sprintf("invalid summoner region: %s", reg.SummonerRegion), "invalid region for League of Legends")
	}

	// Create an URL and ensure it is formatted correctly
	URL := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", reg.SummonerRegion, reg.SummonerName)
	formatURL, err := url.Parse(URL)
	if err != nil {
		return err
	}

	// Use the URL to validate SummonerName against API
	req, err := http.NewRequest(http.MethodGet, formatURL.String(), nil)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	// Set apiKey in header to avoid "403 Unauthorized"
	req.Header.Set("X-Riot-Token", r.apiKey)
	r.mutex.Unlock()

	// Send get-request to API
	resp, err := r.Do(req)
	if err != nil {
		return models.NewAPIErr(err, "Riot")
	}
	defer resp.Body.Close()

	// Ensure that status code is 200 OK, else validation fails
	err = models.AccValStatusCode(resp.StatusCode, "Riot", "invalid username for League of Legends")
	if err != nil {
		return err
	}

	// Decode the response body to get AccountID
	var tmpReg models.SummonerRegistration
	err = json.NewDecoder(resp.Body).Decode(&tmpReg)
	if err != nil {
		return models.NewAPIErr(err, "Riot")
	}

	// Ensure that AccountID is up to date
	reg.AccountID = tmpReg.AccountID
	return nil
}

// UpdateKey updates the riot API key
func (r *Riot) UpdateKey(key string) error {
	// very simple check of the key
	if !strings.HasPrefix(key, "RGAPI-") || len(key) != 42 {
		return models.NewReqErrStr("invalid riot API key", "invalid API key")
	}

	// Create an URL and ensure it is formatted correctly
	URL := "https://EUW1.api.riotgames.com/lol/summoner/v4/summoners/by-name/LOPER"
	formatURL, err := url.Parse(URL)
	if err != nil {
		return err
	}

	// Use the URL to validate SummonerName against API
	req, err := http.NewRequest(http.MethodGet, formatURL.String(), nil)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	// Set apiKey in header to avoid "403 Unauthorized"
	req.Header.Set("X-Riot-Token", r.apiKey)
	r.mutex.Unlock()

	// Send get-request to API
	resp, err := r.Do(req)
	if err != nil {
		return models.NewAPIErr(err, "Riot")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return models.NewReqErrStr("invalid Riot API key", "invalid Riot API key")
	}

	r.mutex.Lock()
	r.apiKey = key
	r.mutex.Unlock()
	return nil
}
