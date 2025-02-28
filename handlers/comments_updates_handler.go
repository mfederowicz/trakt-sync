// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsUpdatesHandler struct for handler
type CommentsUpdatesHandler struct{}

// Handle to handle comments: updates action
func (h CommentsUpdatesHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
