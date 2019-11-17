package models

import (
	"net/http"
)

// Getter is an interface which defines all methods a Getter should provide
// namely, getting a resource based on the given url.
// http.Client is intended to fulfill the interface, but allow for testing without sending requests to another API
//  Similar to client interface, but does not handle auth parameter
type Getter interface {
	Get(url string) (resp *http.Response, err error)
}
