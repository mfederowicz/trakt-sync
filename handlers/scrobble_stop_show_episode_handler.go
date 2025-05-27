// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// ScrobbleStopShowEpisodeHandler struct for handler
type ScrobbleStopShowEpisodeHandler struct{ common CommonLogic }

// Handle to handle scrobble: stop episode type
func (s ScrobbleStopShowEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}

	if options.EpisodeAbs > consts.ZeroValue && len(options.EpisodeCode) > consts.ZeroValue {
		return errors.New("only episode_abs or episode_code at time")
	}

	if len(options.EpisodeCode) > consts.ZeroValue {
		return s.CreateStopScrobbleForEpisodeCode(options, client)
	}

	if options.EpisodeAbs > consts.ZeroValue {
		return s.CreateStopScrobbleForEpisodeAbs(options, client)
	}
	return nil
}

// CreateStopScrobbleForEpisodeCode to handle scrobble: episode code
func (s ScrobbleStopShowEpisodeHandler) CreateStopScrobbleForEpisodeCode(options *str.Options, client *internal.Client) error {
	scrobble, err := s.common.CreateScrobble(client, options)
	if err != nil {
		return fmt.Errorf(consts.ScrobbleError, err)
	}

	result, resp, err := s.common.StopScrobble(client, scrobble, options)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, season:%d, episode:%d", *scrobble.Show.Title, *scrobble.Episode.Season, *scrobble.Episode.Number)
	}

	if err != nil {
		return printer.Errorf(consts.ScrobbleError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, stop episode code scrobble number:%d \n", result.ID)
	}

	return nil
}

// CreateStopScrobbleForEpisodeAbs to handle scrobble: episode abs
func (s ScrobbleStopShowEpisodeHandler) CreateStopScrobbleForEpisodeAbs(options *str.Options, client *internal.Client) error {
	scrobble, err := s.common.CreateScrobble(client, options)
	if err != nil {
		return fmt.Errorf(consts.ScrobbleError, err)
	}

	result, resp, err := s.common.StopScrobble(client, scrobble, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, episode_abs:%d", *scrobble.Show.Title, options.EpisodeAbs)
	}

	if err != nil {
		return printer.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, stop episode abs scrobble number:%d \n", result.ID)
	}

	return nil
}
