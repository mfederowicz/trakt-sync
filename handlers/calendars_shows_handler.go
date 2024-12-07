// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CalendarsShowsHandler struct for handler
type CalendarsShowsHandler struct{}

// Handle to handle calendars: shows action
func (CalendarsShowsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get calendar: " + options.Action)
	result, err := fetchCalendarShows(client, options)
	if err != nil {
		return fmt.Errorf("fetch "+options.Action+" calendar error:%w", err)
	}

	if result == nil {
		return errors.New("empty result")
	}

	printer.Print("Found " + options.Action + " calendar data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")

	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchCalendarShows(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {
	if options.Action == "all-shows" {
		actionType = "all"
	}

	printer.Println("action type:" + actionType)

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetShows(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}
