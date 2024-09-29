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

// WatchlistCmd Returns all items in a user's watchlist filtered by type.
var WatchlistCmd = &Command{
	Name:    "watchlist",
	Usage:   "",
	Summary: "Returns all items in a user's watchlist filtered by type.",
	Help:    `watchlist command`,
}

func watchlistFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	fmt.Println("fetch watchlist lists for:" + options.UserName)

	watchlist, err := fetchWatchlist(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch watchlist error:%w", err)
	}

	if len(watchlist) == consts.ZeroValue {
		return fmt.Errorf("empty watchlist")
	}

	fmt.Printf("Found %d watchlist elements\n", len(watchlist))
	options.Time = cfg.GetOptionTime(options)
	exportJSON := []str.ExportlistItemJSON{}
	findDuplicates := []any{}
	for _, data := range watchlist {
		findDuplicates, exportJSON = cmd.ExportListProcess(data, options, findDuplicates, exportJSON)
	}

	if len(exportJSON) == consts.ZeroValue {
		return fmt.Errorf("warning no data to export, probably a bug")
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
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
	pages, _ := strconv.Atoi(resp.Header.Get(internal.HeaderPaginationPageCount))
	if client.HavePages(page, pages) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := fetchWatchlist(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
