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

	show, err := s.common.FetchShow(client, options)
	if err != nil {
		return fmt.Errorf("fetch show error:%w", err)
	}

	if len(options.EpisodeCode) > consts.ZeroValue {
		return s.CreatePauseScrobbleForEpisodeCode(options, client, show)
	}

	if options.EpisodeAbs > consts.ZeroValue {
		return s.CreatePauseScrobbleForEpisodeAbs(options, client, show)
	}	

	return nil
}

// CreatePauseScrobbleForEpisodeCode to handle pasue scrobble: episode code
func (s ScrobblePauseShowEpisodeHandler) CreatePauseScrobbleForEpisodeCode(options *str.Options, client *internal.Client, show *str.Show) error {
	scrobble := new(str.Scrobble)
	season, number, err := s.common.CheckSeasonNumber(options.EpisodeCode)
	if err != nil {
		return fmt.Errorf("check episode error:%w", err)
	}

	scrobble.Show = new(str.Show)
	scrobble.Show = show
	scrobble.Episode = new(str.Episode)
	scrobble.Episode.Season = season
	scrobble.Episode.Number = number
	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}	

	result, resp, err := s.common.PauseScrobble(client, scrobble)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, season:%d, episode:%d", *show.Title, *season, *number)
	}	

	if err != nil {
		return printer.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, pause episode code scrobble number:%d \n", result.ID)
	}

	return nil
}

// CreatePauseScrobbleForEpisodeAbs to handle pause scrobble: episode abs
func (s ScrobblePauseShowEpisodeHandler) CreatePauseScrobbleForEpisodeAbs(options *str.Options, client *internal.Client, show *str.Show) error {
	scrobble := new(str.Scrobble)
	scrobble.Show = new(str.Show)
	scrobble.Show = show
	scrobble.Episode = new(str.Episode)
	scrobble.Episode.NumberAbs = &options.EpisodeAbs
	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}

	result, resp, err := s.common.PauseScrobble(client, scrobble)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, episode_abs:%d", *show.Title, options.EpisodeAbs)
	}	

	if err != nil {
		return printer.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, pause episode abs scrobble number:%d \n", result.ID)
	}

	return nil
}
