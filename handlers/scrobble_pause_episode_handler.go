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

// ScrobblePauseEpisodeHandler struct for handler
type ScrobblePauseEpisodeHandler struct{ common CommonLogic }

// Handle to handle scrobble: stop movie type
func (s ScrobblePauseEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	episode, _ := s.common.FetchEpisode(client, options)
	scrobble := new(str.Scrobble)
	scrobble.Episode = new(str.Episode)
	scrobble.Episode.IDs = new(str.IDs)
	scrobble.Episode.IDs.Trakt = episode.IDs.Trakt	
	
	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}

	result, resp, err := s.common.PauseScrobble(client, scrobble)
	if err != nil {
		return fmt.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, pause episode scrobble id:%d \n", *result.ID)
	}

	return nil
}
