package main

import (
	"context"
	"net/http"
	"time"

	"github.com/pterm/pterm"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type calEvent struct {
	Name        string
	Date        string
	Description string
	Type        string
	Attachments []*calendar.EventAttachment
	Status      string
	Updated     string
}

// GetCalendar returns calendar events
func GetCalendar(client *http.Client) {
	spinner, _ := pterm.DefaultSpinner.Start("Loading calendar items...")
	spinner.RemoveWhenDone = true

	ctx := context.Background()

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		spinner.Fail("Unable to retrieve Calendar client...")
		Log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339) // Time must be formatted as RFC3339

	// pterm.Info.Println(srv.CalendarList.List().Do())

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(Options.NumItems).OrderBy("startTime").Do()

	if err != nil {
		spinner.Fail("Unable to retrieve calendar events...")
		Log.Fatalf("Unable to retrieve events: %v", err)
	}

	var calEvents []calEvent
	if len(events.Items) == 0 {
		Log.Info("No upcoming events found.")
		pterm.Warning.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := parseDate(item.Start)
			updated := parseUpdated(item.Updated)
			if item.Description == "" {
				item.Description = "n/a"
			}
			calEvents = append(calEvents, calEvent{
				Name:        item.Summary,
				Date:        date,
				Description: item.Description,
				Type:        item.EventType,
				Attachments: item.Attachments,
				Status:      item.Status,
				Updated:     updated,
			})
		}
	}

	spinner.Stop() // Remove spinner

	printEventList(calEvents)
}
