// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinShowEpisodeHandler struct for handler
type CheckinShowEpisodeHandler struct{}

// Handle to handle checkin: episode action
func (h CheckinShowEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
