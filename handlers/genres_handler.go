// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// GenresHandler interface to handle genres
type GenresHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}
