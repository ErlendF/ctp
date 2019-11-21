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
	err error
}


// Do mocks a httpRequest.Do() for testing
func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	// if an error is expected, return it
	if m.err != nil {
		return nil, m.err
	}

	// set http response header
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}

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
	var test = []struct{
		name        string
		payload     *models.SummonerRegistration
		errExpected error
		errHttp     error
	}{
		{"Test OK",&models.SummonerRegistration{SummonerName:"Onijuan",SummonerRegion:"EUW1",AccountID:""},nil,nil},
		{"Test no payload",nil,errors.New("nil summoner registration"),nil},
		//{"Test OK",&models.SummonerRegistration{SummonerName:"",SummonerRegion:"",AccountID:""},errors.New("")},
	} // TODO: need more test cases

	// creating a mockClient item to use the custom "Do" func
	client := &mockClient{}
	riot := New(client, "bigAPIkey")

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			// sets up the client for Do()
			setup := &models.SummonerRegistration{SummonerName:"y",SummonerRegion:"e",AccountID:"s"}
			client.err = tc.errHttp
			client.setup = setup

			// run the function we want to test
			err := riot.ValidateSummoner(tc.payload)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.errExpected, err)
		})
	}
}