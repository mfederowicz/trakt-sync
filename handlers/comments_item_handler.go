// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsItemHandler struct for handler
type CommentsItemHandler struct {}

// Handle to handle comments: comments action
func (h CommentsItemHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
