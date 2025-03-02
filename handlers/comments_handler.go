// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsHandler interface to handle comments module action
type CommentsHandler interface {
	Handle(options *str.Options, client *internal.Client) error 
}
