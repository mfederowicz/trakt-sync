// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
)

// DefaultHandler struct for handler
type DefaultHandler struct{}

// Handle to handle json list item
func (DefaultHandler) Handle(_ *str.Options, _ *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	// movie or show by format imdb
	return findDuplicates, exportJSON, nil
}
