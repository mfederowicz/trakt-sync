// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
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
	var allHandlers = map[string]handlers.Handler{
		"updates":     handlers.PeopleUpdatesHandler{},
		"updated_ids": handlers.PeopleUpdatedIDsHandler{},
		"summary":     handlers.PeopleSummaryHandler{},
		"movies":      handlers.PeopleMoviesHandler{},
		"shows":       handlers.PeopleShowsHandler{},
		"lists":       handlers.PeopleListsHandler{},
		"refresh":     handlers.PeopleRefreshHandler{},
	}
	handler, err := cmd.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.GenActionsUsage([]string{"updates", "updated_ids", "summary", "movies","shows", "lists", "refresh"})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	peopleDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	PeopleCmd.Run = peopleFunc
}
