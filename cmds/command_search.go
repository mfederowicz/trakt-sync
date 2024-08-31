// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

var _searchField str.Slice
var _searchType str.Slice

var (
	_searchAction  = SearchCmd.Flag.String("a", cfg.DefaultConfig().Action, ActionUsage)
	_searchQuery   = SearchCmd.Flag.String("q", cfg.DefaultConfig().Query, QueryUsage)
	_searchID      = SearchCmd.Flag.String("i", cfg.DefaultConfig().ID, IDLookupUsage)
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

func searchFunc(cmd *Command, _ ...string) {

	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	fmt.Println("action:", options.Action)

	switch options.Action {

	case "text-query":

		fmt.Println("Get search: " + options.Action)
		fmt.Printf("search_type: %v\n", options.SearchType.String())
		fmt.Printf("search_field: %v\n", options.SearchField.String())

		result, err := fetchSearchTextQuery(client, options, 1)
		if err != nil {
			fmt.Printf("fetch "+options.Action+" search error:%v", err)
			os.Exit(0)
		}

		if result == nil {
			fmt.Print("empty result")
			os.Exit(0)
		}

		if err == nil {
			if result != nil {
				fmt.Print("Found " + options.Action + " search data \n")
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(result, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No " + options.Action + " search to fetch\n")
			}

		}
	case "id-lookup":

		fmt.Println("Get sarch: " + options.Action)
		fmt.Println("search id_type: " + options.SearchIDType)
		fmt.Println("search id: " + options.ID)
		fmt.Println("search item_type: " + options.SearchType.String())

		result, err := fetchSearchIDLookup(client, options)
		if err != nil {
			fmt.Printf("fetch "+options.Action+" search error:%v", err)
			os.Exit(0)
		}

		if result == nil {
			fmt.Print("empty result")
			os.Exit(0)
		}

		if err == nil {
			if result != nil {
				fmt.Print("Found " + options.Action + " search data \n")
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(result, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No " + options.Action + " search to fetch\n")
			}

		}

	default:
		fmt.Println("possible actions: text-query, id-lookup")
	}

}

var (
	searchDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {

	SearchCmd.Flag.Var(&_searchType, "t", TypeUsage)
	SearchCmd.Flag.Var(&_searchField, "field", FieldUsage)
	SearchCmd.Run = searchFunc
}

func fetchSearchTextQuery(client *internal.Client, options *str.Options, page int) ([]*str.SearchListItem, error) {

	checkRequiredFields(options)
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
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != "" {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > 0 {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + 1
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

	checkRequiredFields(options)
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

func checkRequiredFields(options *str.Options) {

	// Check if the provided module exists in ModuleConfig
	moduleConfig, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		fmt.Printf("Not found config for module '%s'\n", options.Module)
		os.Exit(1)
	}
	// Check search_type flag slice
	if (options.Action == "text-query" && len(options.SearchType) == 0) || !cfg.IsValidConfigTypeSlice(moduleConfig.SearchType, options.SearchType) {
		fmt.Printf("Invalid -t flag values: %v, avaliable values: %v", options.SearchType, moduleConfig.SearchType)
		os.Exit(1)
	}
	// Check search_field flag slice
	if len(options.SearchType) > 0 {
		for _, stype := range options.SearchType {
			if !cfg.IsValidConfigTypeSlice(cfg.SearchFieldConfig[stype], options.SearchField) {
				fmt.Printf("Invalid --field flag values: %v for selected type: %v, avalable values:%v",
					options.SearchField, stype, cfg.SearchFieldConfig[stype])
				os.Exit(1)

			}
		}

	}

	// Check id_type values
	if len(options.SearchIDType) > 0 {
		if !cfg.IsValidConfigType(moduleConfig.SearchIDType, options.SearchIDType) {
			fmt.Printf("Invalid --id_type flag value: %v avalable values:%v",
				options.SearchIDType, moduleConfig.SearchIDType)
			os.Exit(1)

		}
	}

}
