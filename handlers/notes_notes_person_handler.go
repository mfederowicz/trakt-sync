// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesPersonHandler struct for handler
type NotesNotesPersonHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesPersonHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
