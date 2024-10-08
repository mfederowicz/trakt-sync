// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

// ImdbEpisodeHandler struct for handler
type ImdbEpisodeHandler struct{}

// Handle to handle json list item
func (ImdbEpisodeHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	// episode export by format imdb
	findDuplicates = append(findDuplicates, *data.Episode.IDs.Imdb)

	if len(*data.Episode.Title) == consts.ZeroValue {
		notitle := consts.NoEpisodeTitle
		data.Episode.Title = &notitle
	}

	if len(*data.Show.Title) == consts.ZeroValue {
		notitle := consts.NoShowTitle
		data.Show.Title = &notitle
	}

	emap := str.ExportlistItemJSON{
		Imdb:  data.Episode.IDs.Imdb,
		Trakt: data.Episode.IDs.Trakt}
	emap.Uptime(options, data)

	emap.Season = &str.Season{Number: data.Episode.Season}
	emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
	emap.Show = &str.Show{Title: data.Show.Title}

	exportJSON = append(exportJSON, emap)
	return findDuplicates, exportJSON, nil
}
