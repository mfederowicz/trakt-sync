// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersHandler interface to handle users module action
type UsersHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
