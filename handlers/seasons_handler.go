// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// SeasonsHandler interface to handle seasons module action
type SeasonsHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
