// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinHandler interface to handle checkin module action
type CheckinHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
