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

var HistoryCmd = &Command{
	Name:    "history",
	Usage:   "",
	Summary: "returns movies and episodes that a user has watched, sorted by most recent.",
	Help:    `history command`,
}

func historyFunc(cmd *Command, args ...string) {

	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	fmt.Println("fetch history lists for:" + options.UserName)

	history_lists, err := fetchHistoryList(client, options, 1)
	if err != nil {
		fmt.Printf("fetch history list error:%v", err)
		os.Exit(0)
	}

	if len(history_lists) == 0 {
		fmt.Print("empty history lists")
		os.Exit(0)
	}

	fmt.Printf("Found %d history elements\n", len(history_lists))
	options.Time = cfg.GetOptionTime(options)
	export_json := []str.ExportlistItemJson{}
	find_duplicates := []any{}
	for _, data := range history_lists {
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
