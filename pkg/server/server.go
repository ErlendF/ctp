package server

import (
	"fmt"
	"net/http"
	"time"

	"ctp/pkg/models"
)

const writeTimeout, readTimeout, idleTimeout = 60, 60, 60 //timeouts set very high as some requests are suprisingly slow, especially using authentication

// New creates a new http server
func New(port int, apiVer string, organizer models.Organizer) *http.Server {
	handler := newHandler(organizer)
	router := newRouter(handler, apiVer)

	return &http.Server{
		Addr: fmt.Sprintf(":%d", port),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * writeTimeout,
		ReadTimeout:  time.Second * readTimeout,
		IdleTimeout:  time.Second * idleTimeout,
		Handler:      router, // Passing mux router as handler
	}
}
