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
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesJustwatchLinksHandler struct for handler
type MoviesJustwatchLinksHandler struct{}

// Handle to handle movies: justwatch_links action
func (MoviesJustwatchLinksHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := validIDCountryOptions(options, consts.EmptyMovieIDMsg); err != nil {
		return err
	}

	printer.Println("Returns JustWatch links for a movie in the requested country (limited access).")
	result, resp, err := client.Movies.GetMovieJustwatchLinks(cli.ContextFromOptions(options), &options.InternalID, &options.Country)
	if err = watchNowError(consts.JustwatchLinks, consts.Movie, options.InternalID, resp, err); err != nil {
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
