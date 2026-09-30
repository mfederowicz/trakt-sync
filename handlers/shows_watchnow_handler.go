// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsWatchNowHandler struct for handler
type ShowsWatchNowHandler struct{}

// Handle to handle shows: watchnow action
func (ShowsWatchNowHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := validIDCountryOptions(options, consts.EmptyShowIDMsg); err != nil {
		return err
	}

	printer.Println("Returns streaming and watch now sources for a show in the requested country (limited access).")
	opts := uri.ListOptions{Extended: options.ExtendedInfo, Links: options.Links}
	result, resp, err := client.Shows.GetShowWatchNow(cli.ContextFromOptions(options), &options.InternalID, &options.Country, &opts)
	if err = watchNowError(consts.WatchNow, consts.Show, options.InternalID, resp, err); err != nil {
		return err
	}

	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal watchnow error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
