// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

var (
	actionType = "my"
)

// CalendarsHandler interface to handle calendars module action
type CalendarsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// calendarTarget returns the calendar target (my or all) encoded in a {my,all}-* action.
func calendarTarget(action string) string {
	if strings.HasPrefix(action, consts.ActionTypeAll+"-") {
		return consts.ActionTypeAll
	}
	return consts.ActionTypeMy
}

// writeCalendarList writes a fetched calendar list to the output file.
func writeCalendarList(options *str.Options, result []*str.CalendarList) error {
	if result == nil {
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
