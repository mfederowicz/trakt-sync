// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// NotesNotesHandler struct for handler
type NotesNotesHandler struct{ common CommonLogic }

// Handle to handle checkin: checkin action
func (n NotesNotesHandler) Handle(options *str.Options, client *trakt.Client) error {
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

	validTypes := []string{consts.Movie, consts.Show, consts.Season, consts.Episode, consts.Person, consts.History, consts.Collection, consts.Rating}
	if err != nil {
		n.common.GenActionTypeUsage(options, validTypes)
		return unknownValueError("type", options.Type)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(options.Type+":%s", err)
	}

	return nil
}
