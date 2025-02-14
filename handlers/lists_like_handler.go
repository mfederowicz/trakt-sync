// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsLikeHandler struct for handler
type ListsLikeHandler struct{}

// Handle to handle lists: like action
func (h ListsLikeHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
