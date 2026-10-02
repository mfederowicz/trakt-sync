// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncAddToHistoryHandler struct for handler
type SyncAddToHistoryHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_history action
func (m SyncAddToHistoryHandler) Handle(options *str.Options, client *trakt.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("clean history")
	// items hold one entry per play; the cleanup needs every item only once
	toRemove := m.common.CreateItemsToRemove(items.Uniq())

	result, err := m.syncRemoveFromHistory(client, options, &toRemove)
	if err != nil {
		return fmt.Errorf("clean history error:%w", err)
	}

	// the cleanup result keeps its own file; options.Output (-o) is for the add result
	cleanup := *options
	cleanup.Output = fmt.Sprintf(consts.DefaultResultsFormat, options.Module, consts.RemoveFromHistory)

	printer.Println("write cleanup result to:" + cleanup.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	writer.WriteJSON(&cleanup, jsonData)
	waitPageDelay()

	printer.Println("add to history")

	toHistory := m.common.CreateItemsToAdd(items)

	addResult, err := m.syncAddToHistory(client, options, &toHistory)
	if err != nil {
		return fmt.Errorf("add to history error:%w", err)
	}

	jsonDataResult, err := json.MarshalIndent(addResult, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write cleanup result to:" + options.Output)
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (SyncAddToHistoryHandler) syncRemoveFromHistory(client *trakt.Client, options *str.Options, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	result, _, err := client.Sync.RemoveItemsFromHistory(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (SyncAddToHistoryHandler) syncAddToHistory(client *trakt.Client, options *str.Options, items *str.HistoryItems) (*str.AddResult, error) {
	result, _, err := client.Sync.AddItemsToHistory(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
