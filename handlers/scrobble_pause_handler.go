// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// ScrobblePauseHandler struct for handler
type ScrobblePauseHandler struct{ common CommonLogic }

// Handle to handle scrobble: pause action
func (s ScrobblePauseHandler) Handle(options *str.Options, client *trakt.Client) error {
	var handler ScrobbleHandler
	allHandlers := map[string]Handler{
		consts.Movie:       ScrobblePauseMovieHandler{},
		consts.Episode:     ScrobblePauseEpisodeHandler{},
		consts.ShowEpisode: ScrobblePauseShowEpisodeHandler{},
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
