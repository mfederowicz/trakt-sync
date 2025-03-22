// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
)

var (
	_moviesAction        = MoviesCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_moviesMovieIDAction = MoviesCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.MovieIDUsage)
	_moviesPeriod        = MoviesCmd.Flag.String("period", cfg.DefaultConfig().MoviesPeriod, consts.MoviesPeriodUsage)
	_moviesStartDate     = MoviesCmd.Flag.String("start_date", "", consts.StartDateUsage)
)

// MoviesCmd returns movies and episodes that a user has watched, sorted by most recent.
var MoviesCmd = &Command{
	Name:    "movies",
	Usage:   "",
	Summary: "Returns data about movies: trending, popular, list, likes, like, items, comments.",
	Help:    `movies command`,
}

func moviesFunc(cmd *Command, _ ...string) error {	
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	err := cmd.ValidPeriodForModule(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	var handler handlers.MoviesHandler
	
	switch options.Action {
	case "trending":
		handler = handlers.MoviesTrendingHandler{}
	case "popular":
		handler = handlers.MoviesPopularHandler{}
	case "favorited":
		handler = handlers.MoviesFavoritedHandler{}
	case "played":
		handler = handlers.MoviesPlayedHandler{}
	case "watched":
		handler = handlers.MoviesWatchedHandler{}
	case "collected":
		handler = handlers.MoviesCollectedHandler{}
	case "anticipated":
		handler = handlers.MoviesAnticipatedHandler{}
	case "boxoffice":
		handler = handlers.MoviesBoxofficeHandler{}
	case "updates":
		handler = handlers.MoviesUpdatesHandler{}
	case "updated_ids":
		handler = handlers.MoviesUpdatedIDsHandler{}
	case "summary":
		handler = handlers.MoviesSummaryHandler{}
	case "aliases":
		handler = handlers.MoviesAliasesHandler{}

	default:
		printer.Println("possible actions: trending, popular, favorited, played, watched, collected,")
		printer.Println("anticipated, boxoffice, updated, updated_ids,summary,aliases,releases,")
		printer.Println("translations,comments,lists,people,ratings,releated,stats,studios,watiching,videos,refresh")
	}
	
	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	moviesDumpTemplate = ``
)

func init() {
	MoviesCmd.Run = moviesFunc
}
