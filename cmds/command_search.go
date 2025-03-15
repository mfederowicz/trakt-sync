// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

var _searchField str.Slice
var _searchType str.Slice

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

	switch options.Action {
	case "text-query":
		err := runTextQuery(options, client)
		if err != nil {
			return err
		}
	case "id-lookup":
		err := runIDLookup(options, client)
		if err != nil {
			return err
		}

	default:
		printer.Println("possible actions: text-query, id-lookup")
	}
	return nil
}

func runIDLookup(options *str.Options, client *internal.Client) error {
	printer.Println("Get sarch: " + options.Action)
	printer.Println("search id_type: " + options.SearchIDType)
	printer.Println("search id: " + options.ID)
	printer.Println("search item_type: " + options.SearchType.String())

	result, err := fetchSearchIDLookup(client, options)
	if err != nil {
		return fmt.Errorf("fetch "+options.Action+" search error:%w", err)
	}

	if result == nil {
		return errors.New("empty result")
	}

	printer.Print("Found " + options.Action + " search data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")

	writer.WriteJSON(options, jsonData)
	return nil
}

func runTextQuery(options *str.Options, client *internal.Client) error {
	printer.Println("Get search: " + options.Action)
	printer.Printf("search_type: %v\n", options.SearchType.String())
	printer.Printf("search_field: %v\n", options.SearchField.String())
	printer.Println("search id_type: " + options.SearchIDType)

	result, err := fetchSearchTextQuery(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch "+options.Action+" search error:%s", err)
	}

	if result == nil {
		return errors.New("empty result")
	}
	printer.Print("Found " + options.Action + " search data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

var (
	searchDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	SearchCmd.Flag.Var(&_searchType, "t", consts.TypeUsage)
	SearchCmd.Flag.Var(&_searchField, "field", consts.FieldUsage)
	SearchCmd.Run = searchFunc
}

func fetchSearchTextQuery(client *internal.Client, options *str.Options, page int) ([]*str.SearchListItem, error) {
	err := checkRequiredFields(options)
	if err != nil {
		return nil, err
	}

	searchType := options.SearchType.String()
	searchField := options.SearchField.String()
	opts := uri.ListOptions{
		Page:     page,
		Limit:    options.PerPage,
		Extended: options.ExtendedInfo,
		Query:    options.Query,
		Field:    searchField}
	list, resp, err := client.Search.GetTextQueryResults(
		context.Background(),
		&searchType,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := fetchSearchTextQuery(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

func fetchSearchIDLookup(client *internal.Client, options *str.Options) ([]*str.SearchListItem, error) {
	err := checkRequiredFields(options)

	if err != nil {
		return nil, err
	}

	searchType := options.SearchType.String()

	opts := uri.ListOptions{Extended: options.ExtendedInfo, Type: searchType}
	list, _, err := client.Search.GetIDLookupResults(
		context.Background(),
		&options.SearchIDType,
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func noSearchTypeOrInvalidConfigTypeSlice(options *str.Options, slice []string) bool {
	return (options.Action == "text-query" && len(options.SearchType) == consts.ZeroValue) || !cfg.IsValidConfigTypeSlice(slice, options.SearchType)
}

func validSearchIDTypes(options *str.Options, slice []string) bool {
	return len(options.SearchIDType) > consts.ZeroValue && !cfg.IsValidConfigType(slice, options.SearchIDType)
}

func checkSearchFieldFlag(options *str.Options) error {
	if len(options.SearchType) > consts.ZeroValue {
		for _, stype := range options.SearchType {
			if !cfg.IsValidConfigTypeSlice(cfg.SearchFieldConfig[stype], options.SearchField) {
				return fmt.Errorf("invalid --field flag values: %v for selected type: %v, avalable values:%v",
					options.SearchField, stype, cfg.SearchFieldConfig[stype])
			}
		}
	}

	return nil
}

func checkRequiredFields(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}

	// Check search_type flag slice
	if noSearchTypeOrInvalidConfigTypeSlice(options, moduleConfig.SearchType) {
		return fmt.Errorf("invalid -t flag values: %v, avaliable values: %v", options.SearchType, moduleConfig.SearchType)
	}

	// Check search_field flag slice
	sfError := checkSearchFieldFlag(options)
	if sfError != nil {
		return fmt.Errorf("field error:%s", sfError)
	}

	// Check id_type values
	if validSearchIDTypes(options, moduleConfig.SearchIDType) {
		return fmt.Errorf("invalid --id_type flag value: %v avalable values:%v", options.SearchIDType, moduleConfig.SearchIDType)
	}

	return nil
}
