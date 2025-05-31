// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// EpisodesHandler interface to handle episodes module action
type EpisodesHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
