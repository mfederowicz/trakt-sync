// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CheckinHandler interface to handle checkin module action
type CheckinHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
