// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CalendarsMoviesHandler struct for handler
type CalendarsMoviesHandler struct{}

// Handle to handle calendars: movies action
func (CalendarsMoviesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get calendar: " + options.Action + " movies")
	result, err := fetchCalendarMovies(client, options)
	if err != nil {
		return fmt.Errorf("fetch calendar "+options.Action+" error:%w", err)
	}

	if result == nil {
		return fmt.Errorf(consts.EmptyResult)
	}

	printer.Print("Found " + options.Action + " calendar data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchCalendarMovies(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {
	if options.Action == "all-movies" {
		actionType = consts.ActionTypeAll
	}

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetMovies(
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

