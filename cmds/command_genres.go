// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
)

// GenresCmd create or delete active checkins.
var GenresCmd = &Command{
	Name:    "genres",
	Usage:   "",
	Summary: "Get a list of all genres, including names and slugs.",
	Help:    `genres command`,
}

func genresFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.GenresHandler
	switch options.Type {
	case "movies", "shows":
		handler = handlers.GenresTypesHandler{}
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
	genresDumpTemplate = ``
)

func init() {
	GenresCmd.Run = genresFunc
}
