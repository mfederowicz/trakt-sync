// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
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

	fmt.Println("fetch collection lists for:" + options.UserName)

	collection, err := fetchCollectionList(client, options)
	if err != nil {
		return fmt.Errorf("fetch collection error:%w", err)
	}

	if len(collection) == consts.ZeroValue {
		return fmt.Errorf("empty collection")
	}

	fmt.Printf("Found %d collection elements\n", len(collection))
	options.Time = cfg.GetOptionTime(options)
	exportJSON := []str.ExportlistItemJSON{}
	findDuplicates := []any{}
	for _, data := range collection {
		findDuplicates, exportJSON, err = cmd.ExportListProcess(data, options, findDuplicates, exportJSON)
		if err != nil {
			return fmt.Errorf("collection error")
		}

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
	collectionDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	CollectionCmd.Run = collectionFunc
}

func fetchCollectionList(client *internal.Client, options *str.Options) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Sync.GetCollection(
		context.Background(),
		&options.Type,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}
