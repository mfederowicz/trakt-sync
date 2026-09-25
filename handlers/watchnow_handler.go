// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// WatchNowHandler interface to handle watchnow module action
type WatchNowHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
