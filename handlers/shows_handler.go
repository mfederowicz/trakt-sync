// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ShowsHandler interface to handle shows module action
type ShowsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
