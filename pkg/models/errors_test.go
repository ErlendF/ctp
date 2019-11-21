package models

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckStatusCode(t *testing.T) {
	var cases = []struct {
		name     string
		expected error
		code     int
		api      string
		resp     string
	}{
		{"Test OK", nil, http.StatusOK, "test", "oopsie"},
		{"Test unexpected status code", &ExternalAPIError{Err: errors.New("non 200 statuscode"), API: "test", Code: 0}, 0, "test", "oopsie"},
		{"Test not found", &RequestError{Err: fmt.Errorf("invalid response from %s API: %w", "test", ErrNotFound), Response: "oopsie"},
			http.StatusNotFound, "test", "oopsie"},
		{"Test bad request", &RequestError{Err: fmt.Errorf("non 200 statuscode from external API: %s (%d)", "test", http.StatusBadRequest),
			Response: "oopsie"}, http.StatusBadRequest, "test", "oopsie"},
		{"Test unauthorized request", &ExternalAPIError{Err: errors.New("unautorized request to external API"),
			API: "test", Code: http.StatusForbidden}, http.StatusForbidden, "test", "oopsie"},
	}

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckStatusCode(tc.code, tc.api, tc.resp)
			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestAccValStatusCode(t *testing.T) {
	var cases = []struct {
		name     string
		expected error
		code     int
		api      string
		resp     string
	}{
		{"Test OK", nil, http.StatusOK, "test", "oopsie"},
		{"Test forbidden", &ExternalAPIError{Err: errors.New("unautorized request to external API"),
			API: "test", Code: http.StatusForbidden}, http.StatusForbidden, "test", "oopsie"},
		{"Test unexpected status code", &RequestError{Err: fmt.Errorf("non 200 statuscode from external API: %s (%d)", "test", 0),
			Response: "oopsie"}, 0, "test", "oopsie"},
	}

	// tc - test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := AccValStatusCode(tc.code, tc.api, tc.resp)
			assert.Equal(t, tc.expected, err)
		})
	}
}

// this is a very basic test to get more test coverage, because that is a very good metric!
func TestNewErrorFuncs(t *testing.T) {
	testErr := errors.New("test")

	reqErr := &RequestError{Err: testErr, Response: "test response"}
	assert.Equal(t, reqErr, NewReqErrStr("test", "test response"))
	assert.Equal(t, reqErr, NewReqErr(testErr, "test response"))

	unwrappedReqErr := reqErr.Unwrap()
	assert.Equal(t, testErr, unwrappedReqErr)

	reqErrStr := reqErr.Error()
	assert.Equal(t, "test: test response", reqErrStr)

	apiErr := &ExternalAPIError{Err: testErr, API: "test API"}
	assert.Equal(t, apiErr, NewAPIErr(testErr, "test API"))

	unwrappedAPIErr := reqErr.Unwrap()
	assert.Equal(t, testErr, unwrappedAPIErr)

	apiErrResp := apiErr.Respond()
	assert.Equal(t, "Error contacting test API API", apiErrResp)

	apiErrStr := apiErr.Error()
	assert.Equal(t, "error contacting external API test API: test", apiErrStr)
}
