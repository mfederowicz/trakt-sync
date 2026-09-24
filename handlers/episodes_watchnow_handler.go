// Package handlers used to handle module actions
package handlers

import (
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

// EpisodesWatchNowHandler struct for handler
type EpisodesWatchNowHandler struct{}

// Handle to handle episodes: watchnow action
func (EpisodesWatchNowHandler) Handle(options *str.Options, client *internal.Client) error {
	byID := len(options.ID) > consts.ZeroValue
	if byID && len(options.Country) == consts.ZeroValue {
		return errors.New(consts.EmptyCountryMsg)
	}
	if !byID {
		if err := validIDCountryOptions(options, consts.EmptyShowIDMsg); err != nil {
			return err
		}
	}

	printer.Println("Returns streaming and watch now sources for an episode in the requested country (limited access).")
	ctx := client.BuildCtxFromOptions(options)
	opts := uri.ListOptions{Extended: options.ExtendedInfo, Links: options.Links}
	var result map[string]*str.WatchNowSources
	var resp *str.Response
	var err error
	id := options.InternalID
	if byID {
		id = options.ID
		result, resp, err = client.Episodes.GetEpisodeWatchNow(ctx, &options.ID, &options.Country, &opts)
	} else {
		result, resp, err = client.Shows.GetEpisodeWatchNow(ctx, &options.InternalID, &options.Season, &options.Episode, &options.Country, &opts)
	}
	if err = watchNowError(consts.WatchNow, consts.Episode, id, resp, err); err != nil {
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
