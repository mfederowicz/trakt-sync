// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

var (
	_calAction    = CalendarsCmd.Flag.String("a", cfg.DefaultConfig().Action, ActionUsage)
	_calStartDate = CalendarsCmd.Flag.String("start_date", time.Now().Format("2006-01-02"), StartDateUsage)
	_calDays      = CalendarsCmd.Flag.Int("days", 7, DaysUsage)
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

	fmt.Println("action:", options.Action)
	fmt.Println("start_date:", options.StartDate)
	fmt.Println("days:", options.Days)
	switch options.Action {

	case "my-shows", "all-shows":

		fmt.Println("Get calendar: " + options.Action)
		result, err := fetchCalendarShows(client, options)
		if err != nil {
			return fmt.Errorf("fetch "+options.Action+" calendar error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found " + options.Action + " calendar data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")

		writer.WriteJSON(options, jsonData)

	case "my-new-shows", "all-new-shows":

		fmt.Println("Get calendar: " + options.Action)
		result, err := fetchCalendarNewShows(client, options)
		if err != nil {
			return fmt.Errorf("fetch "+options.Action+" calendar error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found " + options.Action + " calendar data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")

		writer.WriteJSON(options, jsonData)

	case "my-season-premieres", "all-season-premieres":

		fmt.Println("Get calendar: " + options.Action + " premieres")
		result, err := fetchCalendarSeasonPremieres(client, options)
		if err != nil {
			return fmt.Errorf("fetch calendar "+options.Action+" premieres error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found " + options.Action + " calendar data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		writer.WriteJSON(options, jsonData)

	case "my-finales", "all-finales":

		fmt.Println("Get calendar: " + options.Action + " finales")
		result, err := fetchCalendarFinales(client, options)
		if err != nil {
			return fmt.Errorf("fetch calendar "+options.Action+" finales error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found " + options.Action + " calendar data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		writer.WriteJSON(options, jsonData)

	case "my-movies", "all-movies":

		fmt.Println("Get calendar: " + options.Action + " movies")
		result, err := fetchCalendarMovies(client, options)
		if err != nil {
			return fmt.Errorf("fetch calendar "+options.Action+" error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found " + options.Action + " calendar data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")

		writer.WriteJSON(options, jsonData)

	case "my-dvd", "all-dvd":

		fmt.Println("Get calendar: " + options.Action + " releases")
		result, err := fetchCalendarDvdReleases(client, options)
		if err != nil {
			return fmt.Errorf("fetch calendar "+options.Action+" error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found " + options.Action + " calendar data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")

		writer.WriteJSON(options, jsonData)

	default:
		fmt.Println("possible actions: {my,all}-shows,{my,all}-new-shows,{my,all}-season-premieres,{my,all}-finales,{my,all}-movies,{my,all}-dvd")
	}

	return nil

}

var (
	calendarsDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	CalendarsCmd.Run = calendarsFunc
}

func fetchCalendarShows(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {

	if options.Action == "all-shows" {
		actionType = "all"
	}

	fmt.Println("action type:" + actionType)

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetShows(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil

}

func fetchCalendarNewShows(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {

	if options.Action == "all-new-shows" {
		actionType = "all"
	}
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetNewShows(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil

}
func fetchCalendarSeasonPremieres(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {

	if options.Action == "all-season-premieres" {
		actionType = "all"
	}

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetSeasonPremieres(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil

}

func fetchCalendarFinales(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {

	if options.Action == "all-finales" {
		actionType = "all"
	}

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetFinales(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil

}

func fetchCalendarMovies(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {

	if options.Action == "all-movies" {
		actionType = "all"
	}

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetMovies(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil

}

func fetchCalendarDvdReleases(client *internal.Client, options *str.Options) ([]*str.CalendarList, error) {

	if options.Action == "all-dvd" {
		actionType = "all"
	}

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Calendars.GetDVDReleases(
		context.Background(),
		&actionType,
		&options.StartDate,
		&options.Days,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil

}
