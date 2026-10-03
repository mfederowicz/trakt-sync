// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CalendarsHandler interface to handle calendars module action
type CalendarsHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}

// calendarTarget returns the calendar target (my or all) encoded in a {my,all}-* action.
func calendarTarget(action string) string {
	if strings.HasPrefix(action, consts.ActionTypeAll+consts.ActionSeparator) {
		return consts.ActionTypeAll
	}
	return consts.ActionTypeMy
}

// calendarOptions builds the query options of a calendar route: extended info and the media filters.
// A calendar is bound by its own start date and days, so the start_date and end_date filters are not sent.
func calendarOptions(options *str.Options) (uri.ListOptions, error) {
	if err := checkMediaFilters(options); err != nil {
		return uri.ListOptions{}, err
	}
	filters := mediaFilters(options)
	filters.StartDate = consts.EmptyString
	filters.EndDate = consts.EmptyString
	return uri.ListOptions{Extended: options.ExtendedInfo, Filters: filters}, nil
}

// writeCalendarList writes a fetched calendar list to the output file.
func writeCalendarList(options *str.Options, result []*str.CalendarList) error {
	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Println("Found " + options.Action + " calendar data")
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal calendar %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
