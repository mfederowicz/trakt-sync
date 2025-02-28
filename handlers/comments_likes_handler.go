// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsLikesHandler struct for handler
type CommentsLikesHandler struct{}

// Handle to handle comments: likes action
func (h CommentsLikesHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
