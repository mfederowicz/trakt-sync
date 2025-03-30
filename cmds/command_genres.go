// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/handlers"
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
	allHandlers := map[string]handlers.Handler{
		"movies": handlers.GenresTypesHandler{},
		"shows":  handlers.GenresTypesHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{"movies", "shows"}
	if err != nil {
		cmd.common.GenTypeUsage(cmd.Name, validTypes)
		return nil
	}

	err = handler.Handle(options, client)
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
