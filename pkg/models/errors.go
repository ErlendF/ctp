package models

import (
	"errors"
	"fmt"
	"net/http"
)

// RequestError indicates the user has made an error in their request
type RequestError struct {
	Response string // Response is used to provide feedback to the user about why the request was bad
	Err      error
}

// ExternalAPIError indicates that there were some error contacting an external API
type ExternalAPIError struct {
	API  string // Indicates which API caused the problem
	Code int
	Err  error
}

//CheckStatusCode checks the status code and returns an ExternalAPIError if it is not 200
func CheckStatusCode(code int, api string, clientResp string) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return &RequestError{Err: fmt.Errorf("invalid response from %s API: %w", api, ErrNotFound), Response: clientResp}
	case http.StatusBadRequest:
		return &RequestError{Err: fmt.Errorf("non 200 statuscode from external API: %s (%d)", api, code), Response: clientResp}
	case http.StatusForbidden:
		return &ExternalAPIError{Err: errors.New("unautorized request to external API"), API: api, Code: code}
	}

	return &ExternalAPIError{Err: errors.New("non 200 statuscode"), API: api, Code: code}
}

//AccValStatusCode checks the status code when validating an account and returns appropriate error
func AccValStatusCode(code int, api string, clientResp string) error {
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusForbidden:
		return &ExternalAPIError{Err: errors.New("unautorized request to external API"), API: api, Code: code}
	}

	return &RequestError{Err: fmt.Errorf("non 200 statuscode from external API: %s (%d)", api, code), Response: clientResp}
}

// Specific errors:

// ErrNotFound indicates that a requested resource was not found
var ErrNotFound = errors.New("not found")

// ErrInvalidID indicates the id is invalid and should not be accepted
var ErrInvalidID = errors.New("invalid id")

// ErrInvalidAuthState defines the error returned if the state for the authentication request does not match the state stored in the cookie
var ErrInvalidAuthState = errors.New("invalid authorization state")

// NewReqErrStr returns a new request error with the given error message and response message
func NewReqErrStr(errStr string, response string) *RequestError {
	return &RequestError{Err: errors.New(errStr), Response: response}
}

// NewReqErr returns a new request error with the given error and response message
func NewReqErr(err error, response string) *RequestError {
	return &RequestError{Err: err, Response: response}
}

// NewAPIErr returns a new external PAI error with the given error message and API name
func NewAPIErr(err error, api string) *ExternalAPIError {
	return &ExternalAPIError{Err: err, API: api}
}

// Methods for the error types:

// Unwrap returns the underlying error
func (e *RequestError) Unwrap() error { return e.Err }

func (e *RequestError) Error() string { return e.Err.Error() + ": " + e.Response }

// Respond returns a string suitable to respond to the user
func (e *ExternalAPIError) Respond() string { return fmt.Sprintf("Error contacting %s API", e.API) }

// Unwrap returns the underlying error
func (e *ExternalAPIError) Unwrap() error { return e.Err }

func (e *ExternalAPIError) Error() string {
	return "error contacting external API " + e.API + ": " + e.Err.Error()
}
