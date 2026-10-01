// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CommentsCommentsHandler struct for handler
type CommentsCommentsHandler struct{}

// Handle to handle checkin: checkin action
func (CommentsCommentsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("generate comment:", options.Type)

	var handler CommentsHandler
	switch options.Type {
	case consts.Movie:
		handler = CommentsCommentsMovieHandler{}
	case consts.Show:
		handler = CommentsCommentsShowHandler{}
	case consts.Season:
		handler = CommentsCommentsSeasonHandler{}
	case consts.Episode:
		handler = CommentsCommentsEpisodeHandler{}
	case consts.List:
		handler = CommentsCommentsListHandler{}
	default:
		printer.Println("possible types: movie,show,season,episode,list")
		return unknownValueError("type", options.Type)
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(options.Type+":%s", err)
	}

	return nil
}
