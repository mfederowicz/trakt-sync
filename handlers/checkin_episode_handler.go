// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinEpisodeHandler struct for handler
type CheckinEpisodeHandler struct{}

// Handle to handle checkin: episode action
func (h CheckinEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
