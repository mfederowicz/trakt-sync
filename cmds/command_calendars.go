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
	_calFilters   = newMediaFilterFlagsWithoutDates(&CalendarsCmd.Flag)
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
	allHandlers := map[string]handlers.Handler{
		consts.MyShows:            handlers.CalendarsShowsHandler{},
		consts.AllShows:           handlers.CalendarsShowsHandler{},
		consts.MyNewShows:         handlers.CalendarsNewShowsHandler{},
		consts.AllNewShows:        handlers.CalendarsNewShowsHandler{},
		consts.MySeasonPremieres:  handlers.CalendarsSeasonPremieresHandler{},
		consts.AllSeasonPremieres: handlers.CalendarsSeasonPremieresHandler{},
		consts.MyFinales:          handlers.CalendarsFinalesHandler{},
		consts.AllFinales:         handlers.CalendarsFinalesHandler{},
		consts.MyMovies:           handlers.CalendarsMoviesHandler{},
		consts.AllMovies:          handlers.CalendarsMoviesHandler{},
		consts.MyDvd:              handlers.CalendarsDvdHandler{},
		consts.AllDvd:             handlers.CalendarsDvdHandler{},
		consts.MyMedia:            handlers.CalendarsMediaHandler{},
		consts.AllMedia:           handlers.CalendarsMediaHandler{},
		consts.MyStreaming:        handlers.CalendarsStreamingHandler{},
		consts.AllStreaming:       handlers.CalendarsStreamingHandler{},
		consts.HotReleases:        handlers.CalendarsHotReleasesHandler{},
		consts.HotPremieres:       handlers.CalendarsHotPremieresHandler{},
		consts.HotNewShows:        handlers.CalendarsHotNewShowsHandler{},
		consts.HotFinales:         handlers.CalendarsHotFinalesHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions := []string{
		"{my,all}_shows", "{my,all}_new_shows", "{my,all}_season_premieres", "{my,all}_finales", "{my,all}_movies", "{my,all}_dvd",
		"{my,all}_media", "{my,all}_streaming", consts.HotReleases, consts.HotPremieres, consts.HotNewShows, consts.HotFinales,
	}
	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validActions)
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	CalendarsCmd.Run = calendarsFunc
}
