// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CalendarsNewShowsHandler struct for handler
type CalendarsNewShowsHandler struct{}

// Handle to handle calendars: shows action
func (CalendarsNewShowsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get calendar: " + options.Action)
	result, err := fetchCalendarNewShows(client, options)
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

func fetchCalendarNewShows(client *trakt.Client, options *str.Options) ([]*str.CalendarList, error) {
	if options.Action == consts.AllNewShows {
		actionType = "all"
	}
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetNewShows(
		cli.ContextFromOptions(options),
		actionType,
		options.StartDate,
		options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}
