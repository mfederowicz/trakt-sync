// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
)

// ImdbShowHandler struct for handler
type ImdbShowHandler struct{}

// Handle to handle json list item
func (h ImdbShowHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	// show-specific logic
	findDuplicates = append(findDuplicates, *data.Show.IDs.Imdb)
	emap := str.ExportlistItemJSON{
		Imdb:  data.Show.IDs.Imdb,
		Trakt: data.Show.IDs.Trakt,
		Title: data.Show.Title}
	emap.Uptime(options, data)

	emap.UpdatedAt = data.UpdatedAt
	exportJSON = append(exportJSON, emap)
	return findDuplicates, exportJSON, nil
}
