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

// HistoryCmd returns movies and episodes that a user has watched, sorted by most recent. 
var HistoryCmd = &Command{
	Name:    "history",
	Usage:   "",
	Summary: "returns movies and episodes that a user has watched, sorted by most recent.",
	Help:    `history command`,
}

func historyFunc(cmd *Command, _ ...string) {

	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	fmt.Println("fetch history lists for:" + options.UserName)

	historyLists, err := fetchHistoryList(client, options, 1)
	if err != nil {
		fmt.Printf("fetch history list error:%v", err)
		os.Exit(0)
	}

	if len(historyLists) == 0 {
		fmt.Print("empty history lists")
		os.Exit(0)
	}

	fmt.Printf("Found %d history elements\n", len(historyLists))
	options.Time = cfg.GetOptionTime(options)
	exportJSON := []str.ExportlistItemJSON{}
	findDuplicates := []any{}
	for _, data := range historyLists {
		findDuplicates, exportJSON = cmd.ExportListProcess(data, options, findDuplicates, exportJSON)
	}

	if len(exportJSON) == 0 {
		print("Warning no data to export, probably a bug")
		os.Exit(1)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
	writer.WriteJSON(options, jsonData)

}

var (
	historyDumpTemplate = ``
)

func init() {
	HistoryCmd.Run = historyFunc
}

func fetchHistoryList(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {

	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetWatchedHistory(
		context.Background(),
		&options.Type,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != "" {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + 1
			nextPageItems, err := fetchHistoryList(client, options, nextPage)
			if err != nil {
				return nil, err
			}

			// Append items from the next page to the current page
			list = append(list, nextPageItems...)

		}

	}

	return list, nil

}
