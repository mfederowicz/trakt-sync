// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ScrobbleStopHandler struct for handler
type ScrobbleStopHandler struct{ common CommonLogic }

// Handle to handle scrobble: stop action
func (s ScrobbleStopHandler) Handle(options *str.Options, client *internal.Client) error {
	var handler ScrobbleHandler
	allHandlers := map[string]Handler{
		"movie":        ScrobbleStopMovieHandler{},
		"episode":      ScrobbleStopEpisodeHandler{},
		"show_episode": ScrobbleStopShowEpisodeHandler{},
	}

	handler, err := s.common.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{"movie", "episode", "show_episode"}
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
