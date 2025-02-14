// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsCommentsHandler struct for handler
type ListsCommentsHandler struct{}

// Handle to handle lists: comments action
func (h ListsCommentsHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
