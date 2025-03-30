// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesItemHandler struct for handler
type NotesItemHandler struct{ common CommonLogic }

// Handle to handle notes: item action
func (n NotesItemHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("get attached item")

	return nil
}
