package blizzard

import (
	"bytes"
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockBlizzard struct {
	setup respSetup
}

func (m *mockBlizzard) Get(url string) (*http.Response, error){
	if !strings.Contains(url, "ow-api.com/v1/stats/") {
		return nil, m.setup.err
	}
	resp := &http.Response{StatusCode:http.StatusOK, Header:make(http.Header, 0)}
	if !strings.Contains(url, "Onijuan-2670") {
		resp.StatusCode = http.StatusNotFound
	}
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
	var test = []struct{
		name    string
		payload *models.Overwatch
		err     string
	}{
		{name:"Test OK", payload:&models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, err: ""},
		{name:"Test invalid BattleTag", payload:&models.Overwatch{BattleTag: "Onyoooo-2670", Platform: "pc", Region: "eu"}, err: "invalid Blizzard battle tag, platform or region"},
		{name:"Test invalid region", payload:&models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "pc"}, err: "invalid Overwatch region"},
		{name:"Test invalid platform", payload:&models.Overwatch{BattleTag: "Onijuan-2670", Platform: "eu", Region: "eu"}, err: "invalid Overwatch platform"},
		{name:"Test no payload", payload:nil, err: "no payload to ValidateBattleUser"},
		//{name:"", payload:&models.Overwatch{BattleTag: "", Platform: "", Region: ""}, err: ""},
	}
	                          // TODO: make error messages consts/line up with 1.13
	getter := &mockBlizzard{}
	ow := New(getter)

	for _, v := range test {
		t.Run(v.name, func (t *testing.T){
			err := ow.ValidateBattleUser(v.payload)
			if err != nil {
				if !strings.Contains(err.Error(), v.err) {
					t.Errorf("Not correct error: |%+v| -- expected |%s|", err, v.err)
				}
			} else {
				if len(v.err) > 0 {
					t.Errorf("Expected error: |%s|", v.err)
				}
			}
		})
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
		{"Test 404",&models.Overwatch{BattleTag: "Onijuan2670", Platform:  "pc", Region:    "eu"}, "","",0,errors.New("invalid Blizzard battle tag or region")},
		{"Test unacceptable response",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "","",0,errors.New("no acceptable response from OW-api")},
		{"Test time [1]",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "0","0",0,nil},
		{"Test time [2]",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "1:0","59:0",1,nil},
		{"Test time [4]",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "4:34:3:2","0",0,errors.New("OW API changed the way time is encoded")},
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
			if err == nil {
				if gem.Time != item.fullTime {
					t.Errorf("Wrong time sum... |%s, %s| -> |%d|", item.cTime, item.qTime, item.fullTime)
				}
				if item.expectedError != nil {
					t.Errorf("Got wrong error... |%v| != |%v|", err, item.expectedError)
				}
				return
			}
			if err != item.expectedError {
				if strings.Contains(err.Error(), item.expectedError.Error()) {
					return
				}
				t.Errorf("Big error from test yes... |%v| != |%v|", err, item.expectedError)
				return
			}

		})
	}

}
