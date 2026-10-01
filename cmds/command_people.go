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

	_peopleReason  = PeopleCmd.Flag.String("r", cfg.DefaultConfig().Reason, consts.ReasonUsage)
	_peopleMessage = PeopleCmd.Flag.String("message", cfg.DefaultConfig().Msg, consts.ReportMsgUsage)
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
		consts.Updates:    handlers.PeopleUpdatesHandler{},
		consts.UpdatedIDs: handlers.PeopleUpdatedIDsHandler{},
		consts.Summary:    handlers.PeopleSummaryHandler{},
		consts.Movies:     handlers.PeopleMoviesHandler{},
		consts.Shows:      handlers.PeopleShowsHandler{},
		consts.Lists:      handlers.PeopleListsHandler{},
		consts.Refresh:    handlers.PeopleRefreshHandler{},

		consts.Report: handlers.PeopleReportHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Updates, consts.UpdatedIDs, consts.Summary, consts.Movies, consts.Shows, consts.Lists, consts.Refresh, consts.Report})
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

var (
	peopleDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	PeopleCmd.Run = peopleFunc
}
