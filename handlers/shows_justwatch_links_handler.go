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

// ShowsJustwatchLinksHandler struct for handler
type ShowsJustwatchLinksHandler struct{}

// Handle to handle shows: justwatch_links action
func (ShowsJustwatchLinksHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validIDCountryOptions(options, consts.EmptyShowIDMsg); err != nil {
		return err
	}

	printer.Println("Returns JustWatch links for a show in the requested country (limited access).")
	result, resp, err := client.Shows.GetShowJustwatchLinks(client.BuildCtxFromOptions(options), &options.InternalID, &options.Country)
	if err = watchNowError(consts.JustwatchLinks, consts.Show, options, resp, err); err != nil {
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
