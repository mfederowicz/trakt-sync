// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// MoviesHandler interface to handle movies module action
type MoviesHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
