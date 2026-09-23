// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// CalendarsStreamingHandler struct for handler
type CalendarsStreamingHandler struct{}

// Handle to handle calendars: {my,all}-streaming action
func (CalendarsStreamingHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get calendar: " + options.Action)
	target := calendarTarget(options.Action)
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.Calendars.GetStreamingReleases(client.BuildCtxFromOptions(options), &target, &options.StartDate, &options.Days, &opts)
	if err != nil {
		return fmt.Errorf("fetch calendar %s error: %w", options.Action, err)
	}

	return writeCalendarList(options, result)
}
