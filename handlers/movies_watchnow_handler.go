// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesWatchNowHandler struct for handler
type MoviesWatchNowHandler struct{}

// Handle to handle movies: watchnow action
func (MoviesWatchNowHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validIDCountryOptions(options, consts.EmptyMovieIDMsg); err != nil {
		return err
	}

	printer.Println("Returns streaming and watch now sources for a movie in the requested country (limited access).")
	opts := uri.ListOptions{Extended: options.ExtendedInfo, Links: options.Links}
	result, resp, err := client.Movies.GetMovieWatchNow(client.BuildCtxFromOptions(options), &options.InternalID, &options.Country, &opts)
	if err = watchNowError(consts.WatchNow, consts.Movie, options.InternalID, resp, err); err != nil {
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
