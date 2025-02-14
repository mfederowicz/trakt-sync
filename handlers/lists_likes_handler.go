// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsLikesHandler struct for handler
type ListsLikesHandler struct{}

// Handle to handle lists: likes action
func (h ListsLikesHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
