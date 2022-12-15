package main

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	// ctx := context.Background()

	GetClient(ReadJSONConfig())
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
