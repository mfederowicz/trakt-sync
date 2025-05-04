// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// LanguagesHandler interface to handle languages
type LanguagesHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
