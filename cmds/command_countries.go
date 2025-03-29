// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
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
	allHandlers := map[string]handlers.Handler{
		"movies": handlers.CountriesTypesHandler{},
		"shows":  handlers.CountriesTypesHandler{},
	}

	handler, err := cmd.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{"movies", "shows"}
	if err != nil {
		cmd.GenTypeUsage(validTypes)
		return nil
	}

	err = handler.Handle(options, client)
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
