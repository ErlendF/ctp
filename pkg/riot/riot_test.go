package riot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockClient is used for setting up the test
type mockClient struct {
	setup string
	err error
}

// Do mocks a httpRequest.Do() for testing
func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	// if a wrong url is sent, return default error
	if !strings.Contains(req.URL.Path, "api.riotgames.com/lol/match/v4") {
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

}