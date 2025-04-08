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

// ScrobbleStartMovieHandler struct for handler
type ScrobbleStartMovieHandler struct{ common CommonLogic }

// Handle to handle scrobble: start movie type
func (s ScrobbleStartMovieHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	movie, _, _ := s.common.FetchMovie(client, options)
	scrobble := new(str.Scrobble)
	scrobble.Movie = movie
	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}

	result, resp, err := s.common.StartScrobble(client, scrobble)
	if err != nil {
		return fmt.Errorf("scrobble error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, movie scrobble action:%d \n", result.Action)
	}

	return nil
}
