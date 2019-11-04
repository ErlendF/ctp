package models

import (
	"net/http"
)

// Client is an interface which defines all methods a Client should provide
// namely, getting a resource based on the given url.
// http.Client is intended to fulfil the interface, but allow for testing without sending requests to another API
//  Similar to Getter interface, but handles auth parameter
type Client interface {
	Do(req *http.Request) (resp *http.Response, err error)
}
