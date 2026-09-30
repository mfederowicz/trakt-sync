// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// Handler interface basic handler
type Handler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
