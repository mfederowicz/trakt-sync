// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

var (
	actionType = "my"
)

// CalendarsHandler interface to handle calendars module action
type CalendarsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
