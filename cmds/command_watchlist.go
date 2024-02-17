package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"trakt-sync/cfg"
	"trakt-sync/internal"
	"trakt-sync/str"
	"trakt-sync/uri"
	"trakt-sync/writer"
)

var WatchlistCmd = &Command{
	Name:    "watchlist",
	Usage:   "",
	Summary: "Returns all items in a user's watchlist filtered by type.",
	Help:    `watchlist command`,
}

func watchlistFunc(cmd *Command, args ...string) {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	fmt.Println("fetch watchlist lists for:" + options.UserName)

	watchlist, err := fetchWatchlist(client, options, 1)
	if err != nil {
		fmt.Printf("fetch watchlist error:%v", err)
		os.Exit(0)
	}

	if len(watchlist) == 0 {
		fmt.Print("empty watchlist")
		os.Exit(0)
	}

	fmt.Printf("Found %d watchlist elements\n", len(watchlist))
	options.Time = cfg.GetOptionTime(options)
	export_json := []str.ExportlistItemJson{}
	find_duplicates := []any{}
	for _, data := range watchlist {
		find_duplicates, export_json = cmd.ExportListProcess(data, options, find_duplicates, export_json)
	}

	if len(export_json) == 0 {
		print("Warning no data to export, probably a bug")
		os.Exit(1)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(export_json, "", "  ")
	writer.WriteJson(options, jsonData)

}

var (
	watchlistDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	WatchlistCmd.Run = watchlistFunc
}

func fetchWatchlist(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetWatchlist(
		context.Background(),
		&options.Type,
		&options.Sort,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != "" {

		pagesInt, _ := strconv.Atoi(pages)

		if pagesInt > 0 && page != pagesInt {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + 1
			nextPageItems, err := fetchWatchlist(client, options, nextPage)
			if err != nil {
				return nil, err
			}

			// Append items from the next page to the current page
			list = append(list, nextPageItems...)

		}

	}

	return list, nil

}
