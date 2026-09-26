// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesHandler struct for handler
type NotesNotesHandler struct{ common CommonLogic }

// Handle to handle checkin: checkin action
func (n NotesNotesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("generate note:", options.Type)

	var handler NotesHandler
	allHandlers := map[string]Handler{
		consts.Movie:      NotesNotesMovieHandler{},
		consts.Show:       NotesNotesShowHandler{},
		consts.Season:     NotesNotesSeasonHandler{},
		consts.Episode:    NotesNotesEpisodeHandler{},
		consts.Person:     NotesNotesPersonHandler{},
		consts.History:    NotesNotesHistoryHandler{},
		consts.Collection: NotesNotesCollectionHandler{},
		consts.Rating:     NotesNotesRatingHandler{},
	}

	handler, err := n.common.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{"movie", "show", "season", "episode", "person", "history", "collection", "rating"}
	if err != nil {
		n.common.GenActionTypeUsage(options, validTypes)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(options.Type+":%s", err)
	}

	return nil
}
