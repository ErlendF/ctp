package riot

import (
	"bytes"
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockClient is used for setting up the test
type mockClient struct {
	setup *models.SummonerRegistration
	code  int
	err   error
}

// Do mocks a httpRequest.Do() for testing
func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	// if an error is expected, return it
	if m.err != nil {
		return nil, m.err
	}

	// set http response header
	resp := &http.Response{StatusCode: m.code, Header: make(http.Header)}

	// add preconfigured response body to THIS response
	body, err := json.Marshal(m.setup)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	// successful request
	return resp, nil
}

func TestRiot_ValidateSummoner(t *testing.T) {
	var test = []struct {
		name        string
		payload     *models.SummonerRegistration
		code        int
		errExpected error
		errHTTP     error
	}{
		{"Test OK", &models.SummonerRegistration{SummonerName: "Onijuan", SummonerRegion: "EUW1", AccountID: "123"}, http.StatusOK, nil, nil},
		{"Test no payload", nil, http.StatusOK, errors.New("nil summoner registration"), nil},
		{"Test no response", &models.SummonerRegistration{SummonerName: "Onijuan", SummonerRegion: "EUW1", AccountID: "123"},
			http.StatusOK, &models.ExternalAPIError{API: "Riot", Code: 0, Err: errors.New("error message")}, errors.New("error message")},
		{"Test invalid username", &models.SummonerRegistration{SummonerName: "Onijuan", SummonerRegion: "EUW1", AccountID: "123"},
			http.StatusNotFound, &models.ExternalAPIError{API: "Riot", Code: 0, Err: errors.New("")}, errors.New("")},
		{"Test invalid region", &models.SummonerRegistration{SummonerName: "Onijuan", SummonerRegion: "gottem", AccountID: "123"},
			http.StatusOK, &models.RequestError{Response: "invalid region for League of Legends",
				Err: errors.New("invalid summoner region: gottem")}, errors.New("")},
		// {"Test ",&models.SummonerRegistration{SummonerName:"",SummonerRegion:"",AccountID:""},http.StatusOK,errors.New(""),errors.New("")},
	}

	// creating a mockClient item to use the custom "Do" func
	client := &mockClient{}
	riot := New(client, "bigAPIkey")

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// sets up the client for Do()
			setup := &models.SummonerRegistration{SummonerName: "y", SummonerRegion: "e", AccountID: "s"}
			client.err = tc.errHTTP
			client.code = tc.code
			client.setup = setup

			// run the function we want to test
			err := riot.ValidateSummoner(tc.payload)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.errExpected, err)
		})
	}
}

func TestRiot_GetRiotPlaytime(t *testing.T) {
	var test = []struct {
		name        string
		payload     *models.SummonerRegistration
		code        int
		errExpected error
		errHTTP     error
	}{
		{"Test OK", &models.SummonerRegistration{SummonerName: "Onijuan", SummonerRegion: "EUW1", AccountID: "123"}, http.StatusOK, nil, nil},
		{"Test no payload", nil, http.StatusOK, errors.New("missing summonerinfo"), nil},
		{"Test no response", &models.SummonerRegistration{SummonerName: "Onijuan", SummonerRegion: "EUW1", AccountID: "123"},
			http.StatusOK, &models.ExternalAPIError{API: "Riot", Code: 0, Err: errors.New("error message")}, errors.New("error message")},
		// {"Test ",&models.SummonerRegistration{SummonerName:"",SummonerRegion:"",AccountID:""},http.StatusOK,errors.New(""),errors.New("")},
	}

	// creating a mockClient item to use the custom "Do" func
	client := &mockClient{}
	riot := New(client, "bigAPIkey")

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// sets up the client for Do()
			setup := &models.SummonerRegistration{SummonerName: "y", SummonerRegion: "e", AccountID: "s"}
			client.err = tc.errHTTP
			client.code = tc.code
			client.setup = setup

			// run the function we want to test
			_, err := riot.GetRiotPlaytime(tc.payload)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.errExpected, err)
		})
	}
}
