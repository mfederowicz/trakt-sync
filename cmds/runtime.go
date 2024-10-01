// Package cmds used for commands modules
package cmds

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"

	"github.com/spf13/afero"
)

// Commands is list of all commands
var Commands = []*Command{
	HelpCmd,
	HistoryCmd,
	WatchlistCmd,
	CollectionCmd,
	UsersListItemsCmd,
	PeopleCmd,
	CalendarsCmd,
	SearchCmd,
}
// runtime consts
const (
	FoundOne = 1
	NotFound = 0
)

func runFoundedModule(cmd *Command, fs afero.Fs, client *internal.Client, config *cfg.Config, args []string) {
	err := cmd.Exec(fs, client, config, args)
	if err != nil {
		fmt.Println(err)
	}
}

// ModulesRuntime core function for process commands
func ModulesRuntime(args []string, config *cfg.Config, client *internal.Client, fs afero.Fs) {
	var found []*Command
	sub, args := args[NotFound], args[FoundOne:]
find:
	for _, cmd := range Commands {
		if sub == cmd.Abbrev {
			found = []*Command{cmd}
			break find
		}
		if strings.HasPrefix(cmd.Name, sub) {
			found = append(found, cmd)
		}
	}

	switch cnt := len(found); cnt {
	case FoundOne:
		runFoundedModule(found[0], fs, client, config, args)
	case NotFound:
		fmt.Fprintf(stdout, "error: unknown command %q\n\n", sub)
		flag.Usage()
	default:
		fmt.Fprintf(stdout, "error: non-unique command prefix %q (matched %d commands)\n\n", sub, cnt)
		flag.Usage()
	}
}
