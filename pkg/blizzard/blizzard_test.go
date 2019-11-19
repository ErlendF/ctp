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

// a mock http.Get that complies with the "Getter" interface
// which allows for customized http responses
func (m *mockBlizzard) Get(url string) (*http.Response, error) {
	// if a wrong url is sent, return default error
	if !strings.Contains(url, "ow-api.com/v1/stats/") {
		return nil, m.setup.err
	}

	// set http response header
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}
	if !strings.Contains(url, "Onijuan-2670") {
		// ensure that a valid user is made, mock invalid user if not
		resp.StatusCode = http.StatusNotFound
	}

	// add preconfigured response body to THIS response
	body, err := json.Marshal(m.setup.resp)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	// successful request
	return resp, nil
}

// struct for setting response body in get request
type respSetup struct {
	resp models.BlizzardResp
	err  error
}

func TestBlizzard_ValidateBattleUser(t *testing.T) {
	var test = []struct {
		name    string
		payload *models.Overwatch
		err     string
	}{
		{name: "Test OK", payload: &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, err: ""},
		{name: "Test invalid BattleTag", payload: &models.Overwatch{BattleTag: "Onyoooo-2670", Platform: "pc", Region: "eu"},
			err: "invalid Blizzard battle tag, platform or region"},
		{name: "Test invalid region", payload: &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "pc"},
			err: "invalid Overwatch region"},
		{name: "Test invalid platform", payload: &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "eu", Region: "eu"},
			err: "invalid Overwatch platform"},
		{name: "Test no payload", payload: nil, err: "no payload to ValidateBattleUser"},
		// {name:"", payload:&models.Overwatch{BattleTag: "", Platform: "", Region: ""}, err: ""},
	}

	// creating a mockBlizzard item to use the custom "Get" func
	getter := &mockBlizzard{}
	ow := New(getter)

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// runs the actual function
			err := ow.ValidateBattleUser(tc.payload)

			// if the error we got does not correspond with the expected error, fail test
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Not correct error: |%+v| -- expected |%s|", err, tc.err)
				}
			} else {
				if len(tc.err) > 0 {
					t.Errorf("Expected error: |%s|", tc.err)
				}
			}
		})
	}
}

func TestBlizzard_GetBlizzardPlaytime(t *testing.T) {
	var testcase = []struct {
		name          string
		payload       *models.Overwatch
		cTime         string
		qTime         string
		fullTime      int
		expectedError error
	}{
		{"Test OK", &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, "0:0:0", "0:0:0", 0, nil},
		{"Test unacceptable response", &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, "", "", 0,
			errors.New("no acceptable response from OW-api")},
		{"Test 404", &models.Overwatch{BattleTag: "Onijuan2670", Platform: "pc", Region: "eu"}, "", "", 0,
			errors.New("invalid Blizzard battle tag or region")},
		{"Test time [1]", &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, "0", "0", 0, nil},
		{"Test time [2]", &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, "1:0", "59:0", 1, nil},
		{"Test time [4]", &models.Overwatch{BattleTag: "Onijuan-2670", Platform: "pc", Region: "eu"}, "4:34:3:2", "0", 0,
			errors.New("OW API changed the way time is encoded")},
		// {"Test ",&models.Overwatch{BattleTag: "Onijuan-2670", Platform:  "pc", Region:    "eu"}, "","",0,errors.New("")},
	}

	// creating a mockBlizzard instance to use the custom "Get" func
	getter := &mockBlizzard{}
	ow := New(getter)

	// Run one test for each of the test cases in array above
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			// set up expected response body with information from test-case
			setup := &respSetup{}
			setup.resp.CompetitiveStats.CareerStats.AllHeroes.Game.TimePlayed = tc.cTime
			setup.resp.QuickPlayStats.CareerStats.AllHeroes.Game.TimePlayed = tc.qTime
			getter.setup = *setup

			// runs the function
			gem, err := ow.GetBlizzardPlaytime(tc.payload)

			// if the error we got does not correspond with the expected error, fail test
			if err == nil {
				if tc.expectedError != nil {
					t.Errorf("Got unexpected error: |%v| != |%v|", err, tc.expectedError)
				}
				// if the resulting time total does not match the test case -> error out
				if gem.Time != tc.fullTime {
					t.Errorf("Unexpected total time played: |%s, %s| -> |%d|", tc.cTime, tc.qTime, tc.fullTime)
				}
				return
			}
			if err != tc.expectedError {
				if strings.Contains(err.Error(), tc.expectedError.Error()) {
					// if the errors contain the expected text -> allow to pass
					return
				}
				t.Errorf("Got unexpected error: |%v| != |%v|", err, tc.expectedError)
				return
			}
		})
	}
}
