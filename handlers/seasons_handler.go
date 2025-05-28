// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// SeasonsHandler interface to handle seasons module action
type SeasonsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
