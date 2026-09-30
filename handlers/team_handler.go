// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// TeamHandler interface to handle team module action
type TeamHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
