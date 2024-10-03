// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// PeopleHandler interface to handle people module action
type PeopleHandler interface {
	Handle(options *str.Options, client *internal.Client) error 
}
