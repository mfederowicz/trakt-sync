// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsCommentHandler struct for handler
type CommentsCommentHandler struct {}

// Handle to handle comments: comment action
func (h CommentsCommentHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
