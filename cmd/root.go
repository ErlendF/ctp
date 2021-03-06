/*
Copyright © 2019 BAKEJ erlend.fonnes@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"ctp/pkg/auth"
	"ctp/pkg/blizzard"
	"ctp/pkg/db"
	"ctp/pkg/jagex"
	"ctp/pkg/models"
	"ctp/pkg/riot"
	"ctp/pkg/user"
	"ctp/pkg/valve"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ctp/pkg/server"

	_ "github.com/joho/godotenv/autoload" // importing .env to os.env
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config struct {
	verbose         bool
	jsonFormatter   bool
	shutdownTimeout int
	clientTimeout   int
	port            int
	fbkey           string
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "ctp",
	Short: "Cloud Technologies (IMT2681) Project",
	Long:  `Cloud Technologies (IMT2681) Project`,

	Run: func(cmd *cobra.Command, args []string) {
		setupLog(config.verbose, config.jsonFormatter)
		logrus.Debugf("Startup config: %+v", config)

		// The Client's Transport typically has internal state (cached TCP
		// connections), so Clients should be reused instead of created as
		// needed. Clients are safe for concurrent use by multiple goroutines
		// - https://golang.org/src/net/http/client.go
		timeout := time.Duration(config.clientTimeout) * time.Second
		client := &http.Client{
			Timeout: timeout,
		}

		// Getting required environment variables (injected by github.com/joho/godotenv/autoload)
		riotAPIKey := os.Getenv("RIOT_API_KEY")
		valveAPIKey := os.Getenv("VALVE_API_KEY")
		clientID := os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
		hmacSecret := os.Getenv("HMAC_SECRET")

		if clientID == "" || clientSecret == "" || hmacSecret == "" || riotAPIKey == "" || valveAPIKey == "" {
			logrus.Fatalf("Invalid environment variables")
		}

		// getting domain for oauth callback, defaulting to localhost
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			domain = "localhost"
		}

		// Initializing each of the provider packages
		riot := riot.New(client, riotAPIKey)
		valve := valve.New(client, valveAPIKey)
		blizzard := blizzard.New(client)
		jagex := jagex.New(client)

		// Getting a database instance
		db, err := db.New(config.fbkey)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to get new Database:%s", err)
		}

		ctx := context.Background()
		ctxC, cancelC := context.WithCancel(ctx)
		defer cancelC()

		// getting a new authenticator, which is passed to the usermanager and server
		auth, err := auth.New(ctxC, db, config.port, domain, clientID, clientSecret, hmacSecret)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to get new Authenticator:%s", err)
		}

		// combining each of the provider packages to fulfill the organizer interface, passing to user manager
		var organizer = struct {
			models.Valve
			models.Riot
			models.Blizzard
			models.Jagex
			models.TokenGenerator
		}{valve, riot, blizzard, jagex, auth}

		um := user.New(db, organizer)
		srv := server.New(config.port, um, auth)

		// Making an channel to listen for errors (later blocking until either error or signal is received)
		errChan := make(chan error)

		// Starting server in a go routine to allow for graceful shutdown
		go func() {
			logrus.Infof("Starting server on port %d", config.port)
			if err := srv.ListenAndServe(); err != nil {
				errChan <- err
			}
		}()

		// Attempting to catch quit via SIGINT (Ctrl+C) to shut down gracefully
		// SIGKILL, SIGQUIT or SIGTERM will not be caught.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		// Blocking until signal or error is received
		select {
		case <-c:
			logrus.Infof("Shutting down server due to interrupt")
		case err := <-errChan:
			logrus.WithError(err).Errorf("Shutting down server due to error")
		}

		ctxT, cancelT := context.WithTimeout(ctx, time.Duration(config.shutdownTimeout)*time.Second)
		defer cancelT()

		// Attempting to shut down the server
		if err := srv.Shutdown(ctxT); err != nil {
			logrus.WithError(err).Fatalf("Unable to gracefully shutdown server")
		}

		logrus.Infoln("Finished shutting down")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init is made and used by cobra
func init() {
	// Reads commandline arguments into config
	rootCmd.Flags().IntVarP(&config.shutdownTimeout, "shutdownTimeout", "s", 15, "Sets the timeout (in seconds) for graceful shutdown")

	rootCmd.Flags().IntVarP(&config.clientTimeout, "clientTimeout", "c", 15,
		"Sets the timeout (in seconds) for the http client which makes requests to the external APIs")

	rootCmd.Flags().IntVarP(&config.port, "port", "p", 80, "Sets the port the API should listen to")
	rootCmd.Flags().BoolVarP(&config.verbose, "verbose", "v", false, "Verbose logging")
	rootCmd.Flags().BoolVarP(&config.jsonFormatter, "jsonFormatter", "j", false, "JSON logging format")
	rootCmd.Flags().StringVarP(&config.fbkey, "fbkey", "f", "./fbkey.json", "Path to the firebase key file")
}

// setupLog initializes logrus logger
func setupLog(verbose, jsonFormatter bool) {
	logLevel := logrus.InfoLevel

	if verbose {
		logLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(os.Stdout)

	if jsonFormatter {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
