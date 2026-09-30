// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// EpisodesHandler interface to handle episodes module action
type EpisodesHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
