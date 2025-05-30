// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CountriesHandler interface to handle countries
type CountriesHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
