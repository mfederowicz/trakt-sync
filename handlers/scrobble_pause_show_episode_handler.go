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

// ScrobblePauseShowEpisodeHandler struct for handler
type ScrobblePauseShowEpisodeHandler struct{ common CommonLogic }

// Handle to handle scrobble: pause show episode type
func (s ScrobblePauseShowEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	if options.EpisodeAbs > consts.ZeroValue && len(options.EpisodeCode) > consts.ZeroValue {
		return errors.New("only episode_abs or episode_code at time")
	}

	if len(options.EpisodeCode) > consts.ZeroValue {
		return s.CreatePauseScrobbleForEpisodeCode(options, client)
	}

	if options.EpisodeAbs > consts.ZeroValue {
		return s.CreatePauseScrobbleForEpisodeAbs(options, client)
	}

	return nil
}

// CreatePauseScrobbleForEpisodeCode to handle pasue scrobble: episode code
func (s ScrobblePauseShowEpisodeHandler) CreatePauseScrobbleForEpisodeCode(options *str.Options, client *internal.Client) error {
	scrobble, err := s.common.CreateScrobble(client, options)
	if err != nil {
		return fmt.Errorf(consts.ScrobbleError, err)
	}

	result, resp, err := s.common.PauseScrobble(client, scrobble, options)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, season:%d, episode:%d", *scrobble.Show.Title, *scrobble.Episode.Season, *scrobble.Episode.Number)
	}

	if err != nil {
		return printer.Errorf(consts.ScrobbleError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, pause episode code scrobble number:%d \n", result.ID)
	}

	return nil
}

// CreatePauseScrobbleForEpisodeAbs to handle pause scrobble: episode abs
func (s ScrobblePauseShowEpisodeHandler) CreatePauseScrobbleForEpisodeAbs(options *str.Options, client *internal.Client) error {
	scrobble, err := s.common.CreateScrobble(client, options)
	if err != nil {
		return fmt.Errorf(consts.ScrobbleError, err)
	}

	result, resp, err := s.common.PauseScrobble(client, scrobble, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, episode_abs:%d", *scrobble.Show.Title, options.EpisodeAbs)
	}

	if err != nil {
		return printer.Errorf(consts.ScrobbleError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, pause episode abs scrobble number:%d \n", result.ID)
	}

	return nil
}
