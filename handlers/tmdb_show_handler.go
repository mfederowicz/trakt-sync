// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
)

// TmdbShowHandler struct for handler
type TmdbShowHandler struct{}

// Handle to handle json list item
func (TmdbShowHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	findDuplicates = append(findDuplicates, *data.Show.IDs.Tmdb)
	emap := str.ExportlistItemJSON{
		Tmdb:  data.Show.IDs.Tmdb,
		Trakt: data.Show.IDs.Trakt,
		Title: data.Show.Title}
	emap.Uptime(options, data)
	exportJSON = append(exportJSON, emap)
	return findDuplicates, exportJSON, nil
}
