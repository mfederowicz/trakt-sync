// Package cmds used for commands modules
package cmds

import (
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

var _searchField str.Slice
var _searchType str.Slice

// legacySearchActions maps the old hyphenated action names to the current ones
var legacySearchActions = map[string]string{
	consts.LegacyIDLookup:  consts.IDLookup,
	consts.LegacyTextQuery: consts.TextQuery,
}

var (
	_searchAction = SearchCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_searchQuery  = SearchCmd.Flag.String("q", cfg.DefaultConfig().Query, consts.QueryUsage)
	_searchID     = SearchCmd.Flag.String("i", cfg.DefaultConfig().ID, IDLookupUsage)
	_searchIDType = SearchCmd.Flag.String("id_type", cfg.DefaultConfig().SearchIDType, IDTypeUsage)
)

// Usage strings in module
const (
	SearchActionUsage = "allow to overwrite action in search command"
	IDLookupUsage     = "allow to overwrite id in search lookup"
	IDTypeUsage       = "allow to overwrite id_type in search lookup"
)

// SearchCmd can use queries or ID lookups
var SearchCmd = &Command{
	Name:    "search",
	Usage:   "",
	Summary: "Searches can use queries or ID lookups",
	Help:    `search command: Queries will search text fields like the title and overview. ID lookups are helpful if you have an external ID and want to get the Trakt ID and info. These methods can search for movies, shows, episodes, people, and str.`,
}

func searchFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	printer.Println("action:", options.Action)

	var handler handlers.SearchHandler
	allHandlers := map[string]handlers.Handler{
		consts.TextQuery: handlers.SearchTextQueryHandler{},
		consts.IDLookup:  handlers.SearchIDLookupHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.TextQuery, consts.IDLookup})
		return nil
	}

	return handler.Handle(options, client)
}

var (
	searchDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	SearchCmd.Flag.Var(&_searchType, "t", consts.TypeUsage)
	SearchCmd.Flag.Var(&_searchField, "field", consts.FieldUsage)
	SearchCmd.Run = searchFunc
}

// normalizeSearchAction replaces a deprecated action name with the current one
func normalizeSearchAction(action string) string {
	current, found := legacySearchActions[action]
	if !found {
		return action
	}

	printer.Printf("action %s is deprecated, use %s\n", action, current)
	return current
}
