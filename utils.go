package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var dateLongFormat = "Mon, Jan 02 2006 01:00 AM"

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

func parseDate(item *calendar.EventDateTime) string {
	var date string
	d := item.DateTime
	if d == "" {
		// parse and format date with no timestamp. (All day event)
		fdate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			Log.Errorf("Unable to parse date: %v", err)
		}
		date = fdate.Format("Mon, Jan 02 2006")
	} else if strings.ContainsAny(d, "T") {
		// parse and format with timestamp
		fdate, err := time.Parse(time.RFC3339, d)
		if err != nil {
			Log.Errorf("Unable to parse date with time: %v", err)
		}
		date = fdate.Format(dateLongFormat)
	}
	return date
}

func parseUpdated(date string) string {
	// parse and format with timestamp
	fdate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		Log.Errorf("Unable to parse date with time: %v", err)
	}

	s := time.Since(fdate).Seconds()
	year := s / 31207680.00
	month := s / 2600640.00
	week := s / 604800.00
	day := s / 86400.00

	if year < 1 && month >= 1 {
		return fmt.Sprintf("~%.0f months ago", month)
	} else if year < 1 && month < 1 {
		return fmt.Sprintf("~%.0f weeks ago", week)
	} else if year < 1 && month < 1 && week >= 1 {
		return fmt.Sprintf("~%.0f days ago", day)
	} else {
		return fmt.Sprintf("~%.0f years ago", year)
	}

}
