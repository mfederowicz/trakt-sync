// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
)

// LanguagesCmd create or delete active checkins.
var LanguagesCmd = &Command{
	Name:    "languages",
	Usage:   "",
	Summary: "Get a list of all languages, including names and codes.",
	Help:    `languages command`,
}

func languagesFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.LanguagesHandler
	allHandlers := map[string]handlers.Handler{
		"movies": handlers.LanguagesTypesHandler{},
		"shows":  handlers.LanguagesTypesHandler{},
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
	languagesDumpTemplate = ``
)

func init() {
	LanguagesCmd.Run = languagesFunc
}
