// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
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
		consts.Movies: handlers.GenresTypesHandler{},
		consts.Shows:  handlers.GenresTypesHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Type, allHandlers)

	validTypes := []string{consts.Movies, consts.Shows}
	if err != nil {
		cmd.common.GenTypeUsage(cmd.Name, validTypes)
		return unknownTypeError(cmd.Name, options.Type)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Type, err)
	}

	return nil
}

func init() {
	GenresCmd.Run = genresFunc
}
