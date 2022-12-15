package main

import (
	"log"
	"os"

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

func main() {
	defer Log.Flush()

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
