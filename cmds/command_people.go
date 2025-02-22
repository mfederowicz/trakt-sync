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
	_action    = PeopleCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_startDate = PeopleCmd.Flag.String("start_date", "", consts.StartDateUsage)
	_personID  = PeopleCmd.Flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
)

// PeopleCmd returns all data for selected person.
var PeopleCmd = &Command{
	Name:    "people",
	Usage:   "",
	Summary: "Returns all data for selected person.",
	Help:    `people command`,
}

func peopleFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.PeopleHandler
	switch options.Action {
	case "updates":
		handler = handlers.PeopleUpdatesHandler{}
	case "updated_ids":
		handler = handlers.PeopleUpdatedIDsHandler{}
	case "summary":
		handler = handlers.PeopleSummaryHandler{}
	case "movies":
		handler = handlers.PeopleMoviesHandler{}
	case "shows":
		handler = handlers.PeopleShowsHandler{}
	case "lists":
		handler = handlers.PeopleListsHandler{}
	case "refresh":
		handler = handlers.PeopleRefreshHandler{}
	default:
		printer.Println("possible actions: updates, updated_ids, summary, movies, shows, lists, refresh")
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s",err)
	}

	return nil
}

var (
	peopleDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	PeopleCmd.Run = peopleFunc
}
