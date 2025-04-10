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
	scrobble, err := s.common.CreateScrobble(client, options)
	if err != nil {
		return fmt.Errorf(consts.ScrobbleError, err)
	}

	result, resp, err := s.common.StartScrobble(client, scrobble)
	if err != nil {
		return fmt.Errorf(consts.ScrobbleError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, movie scrobble action:%d \n", result.Action)
	}

	return nil
}
