// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsLikeHandler struct for handler
type CommentsLikeHandler struct{}

// Handle to handle comments: like action
func (h CommentsLikeHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
