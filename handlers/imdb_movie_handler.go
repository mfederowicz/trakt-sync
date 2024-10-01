// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/str"
)

// ImdbMovieHandler struct for handler
type ImdbMovieHandler struct{}

// Handle to handle json list item
func (h ImdbMovieHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	// movie or show by format imdb
	if !data.Movie.IDs.HaveID("Imdb") {
		noImdb := "no-imdb"
		data.Movie.IDs.Imdb = &noImdb
	}

	findDuplicates = append(findDuplicates, *data.Movie.IDs.Imdb)
	emap := str.ExportlistItemJSON{
		Imdb:  data.Movie.IDs.Imdb,
		Trakt: data.Movie.IDs.Trakt,
		Title: data.Movie.Title}
	emap.Uptime(options, data)
	emap.UpdatedAt = data.UpdatedAt
	emap.Year = data.Movie.Year
	emap.Metadata = data.Metadata
	exportJSON = append(exportJSON, emap)

	return findDuplicates, exportJSON, nil
}
