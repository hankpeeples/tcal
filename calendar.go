package main

import (
	"context"
	"fmt"
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
	ETag        string
	Type        string
	Attachments []*calendar.EventAttachment
	HangoutLink string
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
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		spinner.Fail("Unable to retrieve next ten of the user's calendar events...")
		Log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	var calEvents []calEvent
	if len(events.Items) == 0 {
		Log.Info("No upcoming events found.")
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			calEvents = append(calEvents, calEvent{
				Name:        item.Summary,
				Date:        date,
				Description: item.Description,
				ETag:        item.Etag,
				Type:        item.EventType,
				Attachments: item.Attachments,
				HangoutLink: item.HangoutLink,
				Status:      item.Status,
				Updated:     item.Updated,
			})
		}
	}

	spinner.Stop() // Remove spinner

	printEventList(calEvents)
}

func printEventList(list []calEvent) {
	fmt.Printf("%s \t %s \t %s\n", list[0].Name, list[0].Date, list[0].Status)
}
