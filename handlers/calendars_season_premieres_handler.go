// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CalendarsSeasonPremieresHandler struct for handler
type CalendarsSeasonPremieresHandler struct{}

// Handle to handle calendars: season premieres action
func (CalendarsSeasonPremieresHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get calendar: " + options.Action + " premieres")
	result, err := fetchCalendarSeasonPremieres(client, options)
	if err != nil {
		return fmt.Errorf("fetch calendar "+options.Action+" premieres error:%w", err)
	}

	if result == nil {
		return errors.New(consts.EmptyResult)
	}

	printer.Print("Found " + options.Action + " calendar data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchCalendarSeasonPremieres(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {
	if options.Action == "all-season-premieres" {
		actionType = consts.ActionTypeAll
	}

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetSeasonPremieres(
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

