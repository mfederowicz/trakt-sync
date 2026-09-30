// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// ShowsHandler interface to handle shows module action
type ShowsHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
