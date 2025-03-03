// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsTrendingHandler struct for handler
type CommentsTrendingHandler struct{}

// Handle to handle comments: trending action
func (h CommentsTrendingHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
