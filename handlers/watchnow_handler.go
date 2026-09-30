// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// WatchNowHandler interface to handle watchnow module action
type WatchNowHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
