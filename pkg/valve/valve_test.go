package valve

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// used to mock get request
type mockGetter struct {
	resp *steamResp
	err error
}

func (m *mockGetter) Get(url string) (*http.Response, error) {

	// if an error is set, return it
	if m.err != nil {
		return nil, m.err
	}

	// set http response header
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}

	// add preconfigured response body to THIS response
	body, err := json.Marshal(m.resp)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	// successful request
	return resp, nil
}




