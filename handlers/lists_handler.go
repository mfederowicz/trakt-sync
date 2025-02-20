// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsHandler interface to handle lists module action
type ListsHandler interface {
	Handle(options *str.Options, client *internal.Client) error 
}
