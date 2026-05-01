// Package cmds used for commands modules
package cmds

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
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

	printer.Println("fetch watchlist lists for:" + options.UserName)

	watchlist, err := cmd.common.FetchWatchlist(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch watchlist error:%w", err)
	}

	if len(watchlist) == consts.ZeroValue {
		return errors.New("empty watchlist")
	}

	printer.Printf("Found %d watchlist elements\n", len(watchlist))
	options.Time = cfg.GetOptionTime(options)
	exportJSON := []str.ExportlistItemJSON{}
	findDuplicates := []any{}
	for _, data := range watchlist {
		findDuplicates, exportJSON, err = cmd.ExportListProcess(data, options, findDuplicates, exportJSON)
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
	watchlistDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	WatchlistCmd.Run = watchlistFunc
}
