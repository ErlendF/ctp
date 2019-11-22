package valve

import (
	"bytes"
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// used to mock get request
type mockGetter struct {
	setup respSetup
}

// struct for setting response body in get request
type respSetup struct {
	testRes    *testResp
	resp2      *[]byte
	statusCode int
	err        error
}

type testResp struct {
	Response struct {
		ID64    string `json:"steamid"`
		Code    int    `json:"success"`
		Players []struct {
			ID64           string `json:"steamid"`
			VisibilityCode int    `json:"communityvisibilitystate"`
		} `json:"players"`
	} `json:"response"`
}

func (m *mockGetter) Get(url string) (*http.Response, error) {
	// if an error is set, return it
	if m.setup.err != nil {
		return nil, m.setup.err
	}

	// set http response header
	resp := &http.Response{StatusCode: m.setup.statusCode, Header: make(http.Header)}

	// add preconfigured response body to THIS response
	var body = []byte{}
	var err error
	if m.setup.testRes != nil {
		body, err = json.Marshal(m.setup.testRes)
		if err != nil {
			return nil, err
		}
	}
	if m.setup.resp2 != nil {
		body = *m.setup.resp2
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	// successful request
	return resp, nil
}

func TestValve_ValidateValveAccount(t *testing.T) {
	var test = []struct {
		name          string
		username      string
		ID64          string
		codeResp      int
		vCode         int
		respError     error
		expectedError error
		statusCode    int
	}{
		{name: "Test OK", username: "Onijuan", ID64: "7656119arstarst", vCode: 3, codeResp: 1, expectedError: nil, statusCode: http.StatusOK},
		{name: "Test no username", username: "", ID64: "7656119", vCode: 3, codeResp: 1, expectedError: &models.RequestError{
			Response: "invalid steam account", Err: errors.New("invalid steam account")}, statusCode: http.StatusOK},
		{name: "Test failed Get", username: "Onijuan", ID64: "7656119", vCode: 3,
			expectedError: errors.New("test error"), respError: errors.New("test error")},
		{name: "Test 400 not found", username: "Onijuan", ID64: "7656119", vCode: 3, codeResp: 1, expectedError: &models.RequestError{
			Err:      errors.New("non 200 statuscode from external API: Valve (400)"),
			Response: "invalid steam username"}, statusCode: http.StatusBadRequest},
		{name: "Test invalid account", username: "Onijuan", ID64: "7656119", vCode: 3, codeResp: 0, expectedError: &models.RequestError{
			Err: errors.New("invalid steam account"), Response: "invalid steam account"}, statusCode: http.StatusOK},
		{name: "Test invalid prefix", username: "Onijuan", ID64: "7656f96119", vCode: 3, codeResp: 1, expectedError: &models.RequestError{
			Response: "invalid steam account", Err: errors.New("invalid steam account")}, statusCode: http.StatusOK},
		{name: "Test private account", username: "Onijuan", ID64: "7656119", vCode: 0, codeResp: 1, expectedError: &models.RequestError{
			Err: errors.New("private steam account"), Response: "private steam account"}, statusCode: http.StatusOK},
		// {name:"Test ",username:"Onijuan",ID64:"7656119",codeResp:1,vCode:0,
		// expectedError:errors.New(""),respError:errors.New(""),statusCode:http.StatusOK},
	}

	// creating a mockGetter item to use the custom "Get" func
	getter := &mockGetter{}
	valve := New(getter, "123")

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// setting up the Get() resp according per test_case
			setup := &respSetup{err: tc.respError, statusCode: tc.statusCode, testRes: &testResp{}}
			setup.testRes.Response.ID64 = tc.ID64
			setup.testRes.Response.Code = tc.codeResp
			var tmpPlayer = struct {
				ID64           string `json:"steamid"`
				VisibilityCode int    `json:"communityvisibilitystate"`
			}{ID64: tc.ID64, VisibilityCode: tc.vCode}
			setup.testRes.Response.Players = append(setup.testRes.Response.Players, tmpPlayer)
			getter.setup = *setup

			// runs the actual function
			_, err := valve.ValidateValveAccount(tc.username)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestValve_ValidateValveID(t *testing.T) {
	var test = []struct {
		name          string
		ID64          string
		vCode         int
		respError     error
		expectedError error
		statusCode    int
	}{
		{name: "Test OK", ID64: "7656119arstarst", vCode: 3, expectedError: nil, statusCode: http.StatusOK},
		{name: "Test unauthorized", ID64: "7656119arstarst", vCode: 3, expectedError: &models.ExternalAPIError{
			Err: errors.New("unautorized request to external API"), API: "Valve", Code: http.StatusForbidden}, statusCode: http.StatusForbidden},
		{name: "Test http err", ID64: "7656119arstarst", vCode: 3, expectedError: errors.New("test error"),
			respError: errors.New("test error"), statusCode: http.StatusOK},
		{name: "Test invalid ID", ID64: "765arstars6119arstarst", vCode: 3, expectedError: &models.RequestError{
			Err: errors.New("invalid steam id"), Response: "invalid steam id"}, statusCode: http.StatusOK},
	}

	// creating a mockGetter item to use the custom "Get" func
	getter := &mockGetter{}
	valve := New(getter, "123")

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// setting up the Get() resp according per test_case
			setup := &respSetup{err: tc.respError, statusCode: tc.statusCode, testRes: &testResp{}}
			setup.testRes.Response.ID64 = tc.ID64
			var tmpPlayer = struct {
				ID64           string `json:"steamid"`
				VisibilityCode int    `json:"communityvisibilitystate"`
			}{ID64: tc.ID64, VisibilityCode: tc.vCode}
			setup.testRes.Response.Players = append(setup.testRes.Response.Players, tmpPlayer)
			getter.setup = *setup

			// runs the actual function
			err := valve.ValidateValveID(tc.ID64)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestValve_GetValvePlaytime(t *testing.T) {
	var test = []struct {
		name          string
		ID            string
		resp          *models.ValveResp
		respError     error
		expectedError error
		statusCode    int
	}{
		{name: "Test OK", expectedError: nil, statusCode: http.StatusOK},
		{name: "Test getter error", expectedError: errors.New("test error"), respError: errors.New("test error"), statusCode: http.StatusOK},
		{name: "Test not found", expectedError: &models.RequestError{Err: fmt.Errorf("invalid response from Valve API: %w", models.ErrNotFound),
			Response: "invalid steam id"}, statusCode: http.StatusNotFound},
	}

	// creating a mockGetter item to use the custom "Get" func
	getter := &mockGetter{}
	valve := New(getter, "123")

	var games []models.ValveGames

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// creating test data
			err := faker.FakeData(&games)
			require.Nil(t, err)
			games = append(games, models.ValveGames{Name: "Testing Game", PlaytimeForever: 69420})

			// setting up the Get() resp according per test_case
			setup := &respSetup{err: tc.respError, statusCode: tc.statusCode}
			tmpGameObj := &models.ValveResp{Response: models.ValveResponse{GameCount: len(games), Games: games}}
			body, err := json.Marshal(&tmpGameObj)
			require.Nil(t, err)
			setup.resp2 = &body
			getter.setup = *setup

			// runs the actual function
			_, err = valve.GetValvePlaytime(tc.ID)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
