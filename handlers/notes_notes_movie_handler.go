// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesMovieHandler struct for handler
type NotesNotesMovieHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesMovieHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
