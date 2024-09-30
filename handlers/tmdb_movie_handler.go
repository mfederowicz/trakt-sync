// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
)

// TmdbMovieHandler struct for handler
type TmdbMovieHandler struct{}

// Handle to handle json list item
func (h TmdbMovieHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	findDuplicates = append(findDuplicates, *data.Movie.IDs.Tmdb)
	emap := str.ExportlistItemJSON{
		Tmdb:  data.Movie.IDs.Tmdb,
		Trakt: data.Movie.IDs.Trakt,
		Title: data.Movie.Title}
	emap.Uptime(options, data)
	emap.UpdatedAt = data.UpdatedAt
	exportJSON = append(exportJSON, emap)
	return findDuplicates, exportJSON, nil
}
