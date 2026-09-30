// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// ScrobbleStartHandler struct for handler
type ScrobbleStartHandler struct{ common CommonLogic }

// Handle to handle scrobble: start action
func (s ScrobbleStartHandler) Handle(options *str.Options, client *trakt.Client) error {
	var handler ScrobbleHandler
	allHandlers := map[string]Handler{
		consts.Movie:       ScrobbleStartMovieHandler{},
		consts.Episode:     ScrobbleStartEpisodeHandler{},
		consts.ShowEpisode: ScrobbleStartShowEpisodeHandler{},
	}

	handler, err := s.common.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{consts.Movie, consts.Episode, consts.ShowEpisode}
	if err != nil {
		s.common.GenActionTypeUsage(options, validTypes)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(options.Type+":%s", err)
	}

	return nil
}
