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

// MoviesJustwatchLinksHandler struct for handler
type MoviesJustwatchLinksHandler struct{}

// Handle to handle movies: justwatch_links action
func (MoviesJustwatchLinksHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validMovieCountryOptions(options); err != nil {
		return err
	}

	printer.Println("Returns JustWatch links for a movie in the requested country (limited access).")
	result, resp, err := client.Movies.GetMovieJustwatchLinks(client.BuildCtxFromOptions(options), &options.InternalID, &options.Country)
	if err = movieWatchNowError(consts.JustwatchLinks, options, resp, err); err != nil {
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
