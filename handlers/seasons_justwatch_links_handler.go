// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SeasonsJustwatchLinksHandler struct for handler
type SeasonsJustwatchLinksHandler struct{}

// Handle to handle seasons: justwatch_links action
func (SeasonsJustwatchLinksHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validIDCountryOptions(options, consts.EmptyShowIDMsg); err != nil {
		return err
	}

	printer.Println("Returns JustWatch links for a season in the requested country (limited access).")
	result, resp, err := client.Shows.GetSeasonJustwatchLinks(client.BuildCtxFromOptions(options), &options.InternalID, &options.Season, &options.Country)
	if err = watchNowError(consts.JustwatchLinks, consts.Season, options, resp, err); err != nil {
		return err
	}

	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal justwatch links error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
