// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesRatingHandler struct for handler
type NotesNotesRatingHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesRatingHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
