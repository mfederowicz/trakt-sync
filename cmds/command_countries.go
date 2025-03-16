// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
)

// CountriesCmd create or delete active checkins.
var CountriesCmd = &Command{
	Name:    "countries",
	Usage:   "",
	Summary: "Get a list of all countries, including names and codes.",
	Help:    `countries command`,
}

func countriesFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.CountriesHandler
	switch options.Type {
	case "movies", "shows":
		handler = handlers.CountriesTypesHandler{}
	default:
		printer.Println("possible type: movies,shows")
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Type+":%s", err)
	}

	return nil
}

var (
	countriesDumpTemplate = ``
)

func init() {
	CountriesCmd.Run = countriesFunc
}
