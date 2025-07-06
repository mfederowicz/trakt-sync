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

// HistoryCmd returns movies and episodes that a user has watched, sorted by most recent.
var HistoryCmd = &Command{
	Name:    "history",
	Usage:   "",
	Summary: "Returns movies and episodes that a user has watched, sorted by most recent.",
	Help:    `history command`,
}

func historyFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	printer.Println("fetch history lists for:" + options.UserName)

	historyLists, err := cmd.common.FetchHistoryList(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch history list error:%w", err)
	}

	if len(historyLists) == consts.ZeroValue {
		return errors.New("empty history lists")
	}

	printer.Printf("Found %d history elements\n", len(historyLists))
	options.Time = cfg.GetOptionTime(options)
	exportJSON := []str.ExportlistItemJSON{}
	findDuplicates := []any{}
	for _, data := range historyLists {
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
	historyDumpTemplate = ``
)

func init() {
	HistoryCmd.Run = historyFunc
}
