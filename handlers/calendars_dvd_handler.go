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

// CalendarsDvdHandler struct for handler
type CalendarsDvdHandler struct{}

// Handle to handle calendars: dvd action
func (CalendarsDvdHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get calendar: " + options.Action + " releases")
	result, err := fetchCalendarDvdReleases(client, options)
	if err != nil {
		return fmt.Errorf("fetch calendar "+options.Action+" error:%w", err)
	}

	if result == nil {
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

func fetchCalendarDvdReleases(client *trakt.Client, options *str.Options) ([]*str.CalendarList, error) {
	actionType := calendarTarget(options.Action)

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetDVDReleases(
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
