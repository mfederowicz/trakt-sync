// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesHistoryHandler struct for handler
type NotesNotesHistoryHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesHistoryHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
