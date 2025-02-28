// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)


// CommentsRecentHandler struct for handler
type CommentsRecentHandler struct{}

// Handle to handle comments: recent action
func (h CommentsRecentHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
