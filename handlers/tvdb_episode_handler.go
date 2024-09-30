// Package handlers used to handle list items
package handlers

import (
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

type TvdbEpisodeHandler struct{}

func (h TvdbEpisodeHandler) Handle(options *str.Options, data *str.ExportlistItem, findDuplicates []any, exportJSON []str.ExportlistItemJSON) ([]any, []str.ExportlistItemJSON, error) {
	//fmt.Println("episode export by format tvdb")
	findDuplicates = append(findDuplicates, *data.Episode.IDs.Tvdb)

	if len(*data.Episode.Title) == consts.ZeroValue {
		notitle := consts.NoEpisodeTitle
		data.Episode.Title = &notitle
	}

	if len(*data.Show.Title) == consts.ZeroValue {
		notitle := consts.NoShowTitle
		data.Show.Title = &notitle
	}

	emap := str.ExportlistItemJSON{
		Tvdb:  data.Episode.IDs.Tvdb,
		Trakt: data.Episode.IDs.Trakt}
	emap.Uptime(options, data)

	emap.UpdatedAt = data.UpdatedAt

	emap.Season = &str.Season{Number: data.Episode.Season}
	emap.Episode = &str.Episode{Number: data.Episode.Number, Title: data.Episode.Title}
	emap.Show = &str.Show{Title: data.Show.Title}

	exportJSON = append(exportJSON, emap)
	return findDuplicates, exportJSON, nil
}
