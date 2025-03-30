// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNoteHandler struct for handler
type NotesNoteHandler struct{ common CommonLogic }

// Handle to handle checkin: checkin action
func (n NotesNoteHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("get note")
	return nil
}
