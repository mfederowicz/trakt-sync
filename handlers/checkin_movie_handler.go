// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinMovieHandler struct for handler
type CheckinMovieHandler struct{}

// Handle to handle checkin: checkin action
func (h CheckinMovieHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
