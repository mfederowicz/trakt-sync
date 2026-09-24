// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
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
	if err := validMovieCountryOptions(options); err != nil {
		return err
	}

	printer.Println("Returns streaming and watch now sources for a movie in the requested country (limited access).")
	opts := uri.ListOptions{Extended: options.ExtendedInfo, Links: options.Links}
	result, resp, err := client.Movies.GetMovieWatchNow(client.BuildCtxFromOptions(options), &options.InternalID, &options.Country, &opts)
	if err = movieWatchNowError(consts.WatchNow, options, resp, err); err != nil {
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

func validMovieCountryOptions(options *str.Options) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	if len(options.Country) == consts.ZeroValue {
		return errors.New(consts.EmptyCountryMsg)
	}

	return nil
}

func movieWatchNowError(action string, options *str.Options, resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found movie for:%s", options.InternalID)
	}

	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}

	var forbidden *internal.ForbiddenError
	if errors.As(err, &forbidden) {
		return fmt.Errorf(consts.LimitedAccessMsg, action, err)
	}

	if err != nil {
		return fmt.Errorf("fetch %s error: %w", action, err)
	}

	return nil
}
