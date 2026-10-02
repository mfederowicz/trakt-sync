// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CalendarsHotReleasesHandler struct for handler
type CalendarsHotReleasesHandler struct{}

// Handle to handle calendars: hot_releases action
func (CalendarsHotReleasesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get calendar: " + options.Action)
	opts, err := calendarOptions(options)
	if err != nil {
		return err
	}
	result, _, err := client.Calendars.GetHotReleases(cli.ContextFromOptions(options), options.StartDate, options.Days, &opts)
	if err != nil {
		return fmt.Errorf("fetch calendar %s error: %w", options.Action, err)
	}

	return writeCalendarList(options, result)
}
