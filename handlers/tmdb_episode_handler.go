// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

// TmdbEpisodeHandler struct for handler
type TmdbEpisodeHandler struct{}

// Handle to handle json list item
func (h TmdbEpisodeHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	// episode export by format tmdb
	findDuplicates = append(findDuplicates, *data.Episode.IDs.Tmdb)

	if len(*data.Episode.Title) == consts.ZeroValue {
		notitle := consts.NoEpisodeTitle
		data.Episode.Title = &notitle
	}

	if len(*data.Show.Title) == consts.ZeroValue {
		notitle := consts.NoShowTitle
		data.Show.Title = &notitle
	}

	emap := str.ExportlistItemJSON{
		Tmdb:  data.Episode.IDs.Tmdb,
		Trakt: data.Episode.IDs.Trakt}
	emap.Uptime(options, data)

	emap.Season = &str.Season{Number: data.Episode.Season}
	emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
	emap.Show = &str.Show{Title: data.Show.Title}

	exportJSON = append(exportJSON, emap)
	return findDuplicates, exportJSON, nil
}
