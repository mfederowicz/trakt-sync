// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// PeopleHandler interface to handle people module action
type PeopleHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
