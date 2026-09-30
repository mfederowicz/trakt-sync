// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// UsersHandler interface to handle users module action
type UsersHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
