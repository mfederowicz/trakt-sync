// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesCollectionHandler struct for handler
type NotesNotesCollectionHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
