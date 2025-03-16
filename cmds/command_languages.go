// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
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
	switch options.Type {
	case "movies", "shows":
		handler = handlers.LanguagesTypesHandler{}
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
	languagesDumpTemplate = ``
)

func init() {
	LanguagesCmd.Run = languagesFunc
}
