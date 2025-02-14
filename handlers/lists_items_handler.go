// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsItemsHandler struct for handler
type ListsItemsHandler struct{}

// Handle to handle lists: items action
func (h ListsItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
