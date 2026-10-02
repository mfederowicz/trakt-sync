// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CalendarsMediaHandler struct for handler
type CalendarsMediaHandler struct{}

// Handle to handle calendars: {my,all}_media action
func (CalendarsMediaHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get calendar: " + options.Action)
	target := calendarTarget(options.Action)
	opts, err := calendarOptions(options)
	if err != nil {
		return err
	}
	result, _, err := client.Calendars.GetMedia(cli.ContextFromOptions(options), target, options.StartDate, options.Days, &opts)
	if err != nil {
		return fmt.Errorf("fetch calendar %s error: %w", options.Action, err)
	}

	return writeCalendarList(options, result)
}
