// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// SyncHandler interface to handle sync module action
type SyncHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
