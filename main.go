package main

import (
	"errors"
	"log"
	"os"

	garg "github.com/alexflint/go-arg"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Log is a logger that writes to the tcal log file
var Log *slog.Logger

func init() {
	Log = slog.New()
	h := handler.MustFileHandler("./tcal.log", handler.WithLogLevels(slog.AllLevels))

	Log.PushHandler(h)
}

// Options is the cli options for user configuration
var Options struct {
	NumItems int64 `arg:"-n, --numitems" help:"number of calendar events to pull" default:"10"`
}

func main() {
	defer Log.Flush()

	garg.MustParse(&Options)

	// if token file doesn't exist, get auth and create it
	if _, err := os.Stat("./token.json"); errors.Is(err, os.ErrNotExist) {
		GetClient(ReadJSONConfig())
	}
	// When the user is initialy authenticated, the client needs to be re-configured
	// or the refresh token will be invalid

	// reconnect auth with token file config
	client := GetClient(ReadJSONConfig())

	GetCalendar(client)
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
