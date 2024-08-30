package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

var CollectionCmd = &Command{
	Name:    "collection",
	Usage:   "",
	Summary: "Get all collected items in a user's collection.",
	Help:    `collection command`,
}

func collectionFunc(cmd *Command, args ...string) {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	fmt.Println("fetch collection lists for:" + options.UserName)

	collection, err := fetchCollectionList(client, options)
	if err != nil {
		fmt.Printf("fetch collection error:%v", err)
		os.Exit(0)
	}

	if len(collection) == 0 {
		fmt.Print("empty collection")
		os.Exit(0)
	}

	fmt.Printf("Found %d collection elements\n", len(collection))
	options.Time = cfg.GetOptionTime(options)
	export_json := []str.ExportlistItemJson{}
	find_duplicates := []any{}
	for _, data := range collection {
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
