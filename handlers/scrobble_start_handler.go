// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ScrobbleStartHandler struct for handler
type ScrobbleStartHandler struct{ common CommonLogic }

// Handle to handle scrobble: start action
func (s ScrobbleStartHandler) Handle(options *str.Options, client *internal.Client) error {
	var handler ScrobbleHandler
	allHandlers := map[string]Handler{
		"movie":        ScrobbleStartMovieHandler{},
		"episode":      ScrobbleStartEpisodeHandler{},
		"show_episode": ScrobbleStartShowEpisodeHandler{},
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
