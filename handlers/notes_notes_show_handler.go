// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesShowHandler struct for handler
type NotesNotesShowHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesShowHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
