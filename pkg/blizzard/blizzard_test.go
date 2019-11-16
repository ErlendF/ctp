package blizzard

import (
	"bytes"
	"ctp/pkg/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockBlizzard struct {
	setup respSetup
}

func (m *mockBlizzard) Get(url string) (*http.Response, error){
	validUrl := strings.Contains(url, "ow-api.com/v1/stats/")
	if !validUrl {
		return nil, m.setup.err
	}
	resp := &http.Response{StatusCode:http.StatusOK, Header:make(http.Header, 0)}
	body, err := json.Marshal(m.setup.resp)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	return resp, nil
}

type respSetup struct {
	resp             models.BlizzardResp
	status 			 int
	err              error
}

func TestBlizzard_ValidateBattleUser(t *testing.T) {
	var test = map[models.Overwatch]error{ // TODO: make error messages consts/line up with 1.13
		models.Overwatch{BattleTag:"Onijuan-2670", Platform:"pc", Region:"eu",}: nil,
		models.Overwatch{BattleTag:"Onyoo-7567", Platform:"pc", Region:"eu",}: fmt.Errorf("invalid name"),
		models.Overwatch{BattleTag:"Onijuan-2670", Platform:"pc", Region:"pc",}: fmt.Errorf("invalid region"),
		models.Overwatch{BattleTag:"Onijuan-2670", Platform:"eu", Region:"eu",}: fmt.Errorf("invalid platform"),
	}

	getter := &mockBlizzard{}
	ow := New(getter)

	for k, v := range test {
		err := ow.ValidateBattleUser(&k)
		if err != v {
			t.Errorf("Unexpected error: |%+v| -> |%+v|", &k, err)
		}
	}
}

func TestBlizzard_GetBlizzardPlaytime(t *testing.T) {
	var testcase = []struct{
		name          string
		payload       *models.Overwatch
		cTime         string
		qTime         string
		fullTime      int
		expectedError error
	}{
		{"Test OK",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "0:0:0","0:0:0",0,nil},
		//{"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "","",0,fmt.Errorf("")},
		//{"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "","",0,fmt.Errorf("")},
		//{"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "","",0,fmt.Errorf("")},
	} // TODO: add more test-cases

	getter := &mockBlizzard{}
	ow := New(getter)

	for _, item := range testcase {
		t.Run(item.name, func(t *testing.T){

			setup := &respSetup{}
			setup.resp.CompetitiveStats.CareerStats.AllHeroes.Game.TimePlayed = item.cTime
			setup.resp.QuickPlayStats.CareerStats.AllHeroes.Game.TimePlayed = item.qTime
			getter.setup = *setup

			// TODO: add more checks
			gem, err := ow.GetBlizzardPlaytime(item.payload)
			if err != item.expectedError {
				t.Errorf("Big error from test yes... |%v| != |%v|", err, item.expectedError)
				return
			}

			if err == nil {
				if gem.Time != item.fullTime {
					t.Errorf("Wrong time sum... |%s, %s| -> |%d|", item.cTime, item.qTime, item.fullTime)
				}
			}





		})
	}

}
