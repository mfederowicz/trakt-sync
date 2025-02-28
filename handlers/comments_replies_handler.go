// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsRepliesHandler struct for handler
type CommentsRepliesHandler struct{}

// Handle to handle comments: replies action
func (h CommentsRepliesHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
