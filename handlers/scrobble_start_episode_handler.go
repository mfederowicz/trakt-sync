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

// ScrobbleStartEpisodeHandler struct for handler
type ScrobbleStartEpisodeHandler struct{ common CommonLogic }

// Handle to handle scrobble: start episode type
func (s ScrobbleStartEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	
	scrobble, err := s.common.CreateScrobble(client, options)
	if err != nil {
		return fmt.Errorf("scrobble error:%w", err)
	}

	result, resp, err := s.common.StartScrobble(client, scrobble)
	if err != nil {
		return fmt.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode scrobble action:%d \n", result.Action)
	}

	return nil
}
