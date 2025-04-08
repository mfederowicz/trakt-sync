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

// ScrobblePauseMovieHandler struct for handler
type ScrobblePauseMovieHandler struct{ common CommonLogic }

// Handle to handle scrobble: stop movie type
func (s ScrobblePauseMovieHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	movie, _, _ := s.common.FetchMovie(client, options)
	scrobble := new(str.Scrobble)
	scrobble.Movie = movie
	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}

	result, resp, err := s.common.PauseScrobble(client, scrobble)
	if err != nil {
		return fmt.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, pause movie scrobble id:%d \n", *result.ID)
	}

	return nil
}
