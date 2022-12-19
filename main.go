package main

import (
	garg "github.com/alexflint/go-arg"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
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

	if NeedsLogin() {
		GetClient(ReadJSONConfig())
	}

	// reconnect auth with token file config
	client := GetClient(ReadJSONConfig())

	GetCalendar(client)
}
