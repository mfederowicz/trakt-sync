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
	"github.com/mfederowicz/trakt-sync/writer"
)

// CalendarsNewShowsHandler struct for handler
type CalendarsNewShowsHandler struct{}

// Handle to handle calendars: {my,all}_new_shows action
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
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchCalendarNewShows(client *trakt.Client, options *str.Options) ([]*str.CalendarList, error) {
	actionType := calendarTarget(options.Action)
	opts, err := calendarOptions(options)
	if err != nil {
		return nil, err
	}
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
