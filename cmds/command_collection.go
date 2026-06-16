// Package cmds used for commands modules
package cmds

import (
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

// CollectionCmd get all collected items in a user's collection.
var CollectionCmd = &Command{
	Name:    "collection",
	Usage:   "",
	Summary: "Get all collected items in a user's collection.",
	Help:    `collection command`,
}

func collectionFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	printer.Println("fetch collection lists for:" + options.UserName)

	collection, err := fetchCollectionList(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch collection error:%w", err)
	}

	if len(collection) == consts.ZeroValue {
		return errors.New("empty collection")
	}

	printer.Printf("Found %d collection elements\n", len(collection))
	options.Time = cfg.GetOptionTime(options)
	exportJSON := []str.ExportlistItemJSON{}
	findDuplicates := []any{}
	for _, data := range collection {
		findDuplicates, exportJSON, err = cmd.ExportListProcess(data, options, findDuplicates, exportJSON)
		if err != nil {
			return errors.New("collection error")
		}
	}

	if len(exportJSON) == consts.ZeroValue {
		return errors.New("warning no data to export, probably a bug")
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

var (
	collectionDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	CollectionCmd.Run = collectionFunc
}

func fetchCollectionList(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetCollection(
		client.BuildCtxFromOptions(options),
		&options.Type,
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
		nextPageItems, err := fetchCollectionList(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}
