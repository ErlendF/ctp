package valve

import (
	"bytes"
	"ctp/pkg/models"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// used to mock get request
type mockGetter struct {
	setup respSetup
}

// struct for setting response body in get request
type respSetup struct {
	resp       steamResp
	statusCode int
	err        error
}

func (m *mockGetter) Get(url string) (*http.Response, error) {

	// if an error is set, return it
	if m.setup.err != nil {
		return nil, m.setup.err
	}

	// set http response header
	resp := &http.Response{StatusCode: m.setup.statusCode, Header: make(http.Header)}

	// add preconfigured response body to THIS response
	body, err := json.Marshal(m.setup.resp)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	// successful request
	return resp, nil
}

func TestValve_ValidateValveAccount(t *testing.T) {
	var test = []struct{
		name string
		username string
		ID64 string
		codeResp int
		respError error
		expectedError error
		statusCode int
	}{
		{name:"Test OK",username:"Onijuan",ID64:"7656119arstarst",codeResp:1,expectedError:nil,statusCode:http.StatusOK},
		{name:"Test no username",username:"",ID64:"7656119",codeResp:1,expectedError:&models.RequestError{Response:"invalid steam account",Err:errors.New("invalid steam account")},statusCode:http.StatusOK},
		{name:"Test failed Get",username:"Onijuan",expectedError:errors.New("test error"),respError:errors.New("test error")},
		{name:"Test 400 not found",username:"Onijuan",ID64:"7656119",codeResp:1,expectedError:&models.RequestError{Err:errors.New("non 200 statuscode from external API: Valve (400)"),Response:"invalid steam username",},statusCode:http.StatusBadRequest},
		{name:"Test invalid account",username:"Onijuan",ID64:"7656119",codeResp:0,expectedError:&models.RequestError{Err:errors.New("invalid steam account"), Response:"invalid steam account"},statusCode:http.StatusOK},
		{name:"Test invalid prefix",username:"Onijuan",ID64:"7656f96119",codeResp:1,expectedError:&models.RequestError{Response:"invalid steam account",Err:errors.New("invalid steam account")},statusCode:http.StatusOK},
		//{name:"Test ",username:"Onijuan",ID64:"7656119",codeResp:1,expectedError:errors.New(""),respError:errors.New(""),statusCode:http.StatusOK},
	}

	// creating a mockGetter item to use the custom "Get" func
	getter := &mockGetter{}
	valve := New(getter, "123")

	// run a test for each of the test items (array above)
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {

			// setting up the Get() resp according per test_case
			setup := &respSetup{err: tc.respError, statusCode: tc.statusCode}
			setup.resp.Response.ID64 = tc.ID64
			setup.resp.Response.Code = tc.codeResp
			getter.setup = *setup

			// runs the actual function
			_, err := valve.ValidateValveAccount(tc.username)

			// if the error we got does not correspond with the expected error, fail test
			assert.Equal(t, tc.expectedError, err)
		})
	}
}


