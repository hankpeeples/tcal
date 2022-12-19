package main

import (
	"errors"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// NeedsLogin will check if the token file already exists or not
func NeedsLogin() bool {
	// if token file doesn't exist
	if _, err := os.Stat("./token.json"); errors.Is(err, os.ErrNotExist) {
		return true
	}
	// When the user is initialy authenticated, the client needs to be re-configured
	// or the refresh token will be invalid
	return false
}

// ReadJSONConfig reads the json file holding app client secret values
func ReadJSONConfig() *oauth2.Config {
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Printf("Can't read file... %s", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
	}

	return config
}
