package server

import (
	"fmt"
	"net/http"
	"time"

	"ctp/pkg/models"
)

const writeTimeout, readTimeout, idleTimeout = 15, 30, 30

// New creates a new http server
func New(port int, um models.UserManager) *http.Server {
	handler := newHandler(um)
	router := newRouter(handler)

	return &http.Server{
		Addr: fmt.Sprintf(":%d", port),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * writeTimeout,
		ReadTimeout:  time.Second * readTimeout,
		IdleTimeout:  time.Second * idleTimeout,
		Handler:      router, // Passing mux router as handler
	}
}
