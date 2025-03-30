// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesSeasonHandler struct for handler
type NotesNotesSeasonHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesSeasonHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
