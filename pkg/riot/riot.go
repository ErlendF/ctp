package riot

import (
	"ctp/pkg/models"
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
func (r *Riot) GetRiotPlaytime() (*models.Game, error) {
	logrus.Debugf("GetLolPlaytime")

	// db.UpdateGame("League", 9, "117575669351657432712")
	game := &models.Game{Name: "League", Time: 9}
	return game, nil
}

//ValidateSummoner validates the summoner
func (r *Riot) ValidateSummoner(regInfo *models.SummonerRegistration) error {
	var regions = []string{"RU", "KR", "BR1", "OC1", "JP1", "NA1", "EUN1", "EUW1", "TR1", "LA1", "LA2"}
	// var err error

	validRegion := false

	//Validating region
	for _, r := range regions {
		if r == regInfo.SummonerRegion {
			validRegion = true
		}
	}

	if !validRegion {
		return fmt.Errorf("invalid region")
	}

	//Validating name
	URL := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", regInfo.SummonerRegion, regInfo.SummonerName)

	formatURL, err := url.Parse(URL)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, formatURL.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Riot-Token", r.apiKey)

	resp, err := r.Client.Do(req)

	if err != nil{
		return err
	}
	if resp.StatusCode != http.StatusOK  {
		return fmt.Errorf(string(resp.StatusCode))
	}

	return nil
}
