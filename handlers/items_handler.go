// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
)

// ItemsHandler interface to handle list item
type ItemsHandler interface {
	Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error)
}
