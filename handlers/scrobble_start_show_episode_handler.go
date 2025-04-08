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

// ScrobbleStartShowEpisodeHandler struct for handler
type ScrobbleStartShowEpisodeHandler struct{ common CommonLogic }

// Handle to handle scrobble: start show_episode type
func (s ScrobbleStartShowEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
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
		return s.CreateScrobbleForEpisodeCode(options, client, show)
	}

	if options.EpisodeAbs > consts.ZeroValue {
		return s.CreateScrobbleForEpisodeAbs(options, client, show)
	}

	return nil
}

// CreateScrobbleForEpisodeCode to handle scrobble: episode code
func (s ScrobbleStartShowEpisodeHandler) CreateScrobbleForEpisodeCode(options *str.Options, client *internal.Client, show *str.Show) error {
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

	result, resp, err := s.common.StartScrobble(client, scrobble)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, season:%d, episode:%d", *show.Title, *season, *number)
	}	

	if err != nil {
		return printer.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode code scrobble number:%d \n", result.ID)
	}

	return nil
}

// CreateScrobbleForEpisodeAbs to handle scrobble: episode abs
func (s ScrobbleStartShowEpisodeHandler) CreateScrobbleForEpisodeAbs(options *str.Options, client *internal.Client, show *str.Show) error {
	scrobble := new(str.Scrobble)
	scrobble.Show = new(str.Show)
	scrobble.Show = show
	scrobble.Episode = new(str.Episode)
	scrobble.Episode.NumberAbs = &options.EpisodeAbs
	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}

	result, resp, err := s.common.StartScrobble(client, scrobble)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, episode_abs:%d", *show.Title, options.EpisodeAbs)
	}	

	if err != nil {
		return printer.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode abs scrobble number:%d \n", result.ID)
	}

	return nil
}
