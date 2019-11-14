package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//NonOK specifies errormessage returned when statuscode returned from api is not 200
const NonOK = "non 200 statuscode"

//nonOKfmt is used to format the error message returned by CheckStatusCode
const nonOKfmt = NonOK + ": %d"

const trim = "non 200 statuscode: "

//InvalidAuthState defines the error returned if the state for the authenticaiton request does not match the state stored in the cookie
const InvalidAuthState = "Invalid state"

//NotFound defines the error returned if a resource could not be found
const NotFound = "NotFound"

//ClientError indicates the error was caused by the client
const ClientError = "Client error"

//ServerError indicates the error was caused by something on the serverside
const ServerError = "Server error"

//CheckStatusCode checks the status code and returns an error if it is not 200
func CheckStatusCode(code int) error {
	if code != http.StatusOK {
		return fmt.Errorf(nonOKfmt, code)
	}

	return nil
}

//CheckNotFound checks the error to see if it is a 404 Not found error from a request
func CheckNotFound(err error) bool {
	if err == nil {
		return false
	}

	if !strings.Contains(err.Error(), NonOK) {
		return false
	}
	statusCode := strings.TrimPrefix(err.Error(), trim)
	if statusCode == strconv.Itoa(http.StatusNotFound) { // NOT HARD CODED, YES
		return true
	}

	return false
}

//GetHTTPErrorClass returns the class of the HTTP status code, if applicable
func GetHTTPErrorClass(err error) string {
	if err == nil {
		return ""
	}

	if !strings.Contains(err.Error(), NonOK) {
		return ""
	}

	statusCode := strings.TrimPrefix(err.Error(), trim)
	if len(statusCode) < 1 {
		return ""
	}

	class := string(statusCode[0])

	switch class {
	case "4":
		return ClientError
	case "5":
		return ServerError
	}

	return ""
}
