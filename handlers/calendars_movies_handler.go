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

// CalendarsMoviesHandler struct for handler
type CalendarsMoviesHandler struct{}

// Handle to handle calendars: movies action
func (CalendarsMoviesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get calendar: " + options.Action + " movies")
	result, err := fetchCalendarMovies(client, options)
	if err != nil {
		return fmt.Errorf("fetch calendar "+options.Action+" error:%w", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
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

func fetchCalendarMovies(client *trakt.Client, options *str.Options) ([]*str.CalendarList, error) {
	actionType := calendarTarget(options.Action)

	opts, err := calendarOptions(options)
	if err != nil {
		return nil, err
	}
	list, _, err := client.Calendars.GetMovies(
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
