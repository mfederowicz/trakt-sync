// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesEpisodeHandler struct for handler
type NotesNotesEpisodeHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
