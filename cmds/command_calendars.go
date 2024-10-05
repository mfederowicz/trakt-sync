// Package cmds used for commands modules
package cmds

import (
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
)

var (
	_calAction    = CalendarsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_calStartDate = CalendarsCmd.Flag.String("start_date", time.Now().Format("2006-01-02"), consts.StartDateUsage)
	_calDays      = CalendarsCmd.Flag.Int("days", 7, consts.DaysUsage)
	actionType    = "my"
)

// CalendarsCmd process selected user calendars
var CalendarsCmd = &Command{
	Name:    "calendars",
	Usage:   "",
	Summary: "By default, the calendar will return all shows or movies for the specified time period and can be global or user specific.",
	Help:    `calendars command`,
}

func calendarsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	printer.Println("action:", options.Action)
	printer.Println("start_date:", options.StartDate)
	printer.Println("days:", options.Days)
	var handler handlers.CalendarsHandler
	switch options.Action {
	case "my-shows", "all-shows":
		handler = handlers.CalendarsShowsHandler{}
	case "my-new-shows", "all-new-shows":
		handler = handlers.CalendarsNewShowsHandler{}
	case "my-season-premieres", "all-season-premieres":
		handler = handlers.CalendarsSeasonPremieresHandler{}
	case "my-finales", "all-finales":
		handler = handlers.CalendarsFinalesHandler{}
	case "my-movies", "all-movies":
		handler = handlers.CalendarsMoviesHandler{}
	case "my-dvd", "all-dvd":
		handler = handlers.CalendarsDvdHandler{}
	default:
		printer.Println("possible actions: {my,all}-shows,{my,all}-new-shows,{my,all}-season-premieres,{my,all}-finales,{my,all}-movies,{my,all}-dvd")
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	calendarsDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	CalendarsCmd.Run = calendarsFunc
}
