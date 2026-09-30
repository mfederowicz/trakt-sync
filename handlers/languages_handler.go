// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// LanguagesHandler interface to handle languages
type LanguagesHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}
