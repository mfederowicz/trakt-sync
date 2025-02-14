// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsListHandler struct for handler
type ListsListHandler struct{}

// Handle to handle lists: list action
func (h ListsListHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
