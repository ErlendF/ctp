package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//NonOK specifies errormessage returned when statuscode returned from api is not 200
const NonOK = "Non 200 statuscode"

//nonOKfmt is used to format the error message returned by CheckStatusCode
const nonOKfmt = NonOK + ": %d"

const trim = "Non 200 statuscode: "

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
