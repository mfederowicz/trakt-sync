// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
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

	fmt.Println("action:", options.Action)

	switch options.Action {

	case "text-query":

		fmt.Println("Get search: " + options.Action)
		fmt.Printf("search_type: %v\n", options.SearchType.String())
		fmt.Printf("search_field: %v\n", options.SearchField.String())

		result, err := fetchSearchTextQuery(client, options, consts.DefaultPage)
		if err != nil {
			return fmt.Errorf("fetch "+options.Action+" search error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}
		fmt.Print("Found " + options.Action + " search data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JsonDataFormat)

		writer.WriteJSON(options, jsonData)

	case "id-lookup":

		fmt.Println("Get sarch: " + options.Action)
		fmt.Println("search id_type: " + options.SearchIDType)
		fmt.Println("search id: " + options.ID)
		fmt.Println("search item_type: " + options.SearchType.String())

		result, err := fetchSearchIDLookup(client, options)
		if err != nil {
			return fmt.Errorf("fetch "+options.Action+" search error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}
		
		fmt.Print("Found " + options.Action + " search data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")

		writer.WriteJSON(options, jsonData)

	default:
		fmt.Println("possible actions: text-query, id-lookup")
	}
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
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, Query: options.Query, Field: searchField}
	list, resp, err := client.Search.GetTextQueryResults(
		context.Background(),
		&searchType,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != consts.EmptyString {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > consts.ZeroValue {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + consts.NextPageStep
			nextPageItems, err := fetchSearchTextQuery(client, options, nextPage)
			if err != nil {
				return nil, err
			}

			// Append items from the next page to the current page
			list = append(list, nextPageItems...)

		}

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

func checkRequiredFields(options *str.Options) error {

	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	// Check search_type flag slice
	if (options.Action == "text-query" && len(options.SearchType) == consts.ZeroValue) || !cfg.IsValidConfigTypeSlice(moduleConfig.SearchType, options.SearchType) {
		return fmt.Errorf("invalid -t flag values: %v, avaliable values: %v", options.SearchType, moduleConfig.SearchType)
	}
	// Check search_field flag slice
	if len(options.SearchType) > consts.ZeroValue {
		for _, stype := range options.SearchType {
			if !cfg.IsValidConfigTypeSlice(cfg.SearchFieldConfig[stype], options.SearchField) {
				return fmt.Errorf("invalid --field flag values: %v for selected type: %v, avalable values:%v",
					options.SearchField, stype, cfg.SearchFieldConfig[stype])

			}
		}

	}

	// Check id_type values
	if len(options.SearchIDType) > consts.ZeroValue {
		if !cfg.IsValidConfigType(moduleConfig.SearchIDType, options.SearchIDType) {
			return fmt.Errorf("invalid --id_type flag value: %v avalable values:%v",
				options.SearchIDType, moduleConfig.SearchIDType)

		}
	}

	return nil

}
