// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncAddToHistoryHandler struct for handler
type SyncAddToHistoryHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_history action
func (m SyncAddToHistoryHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("clean history")
	toRemove := m.common.CreateItemsToRemove(items)

	result, err := m.syncRemoveFromHistory(client, options, &toRemove)
	if err != nil {
		return fmt.Errorf("clean history error:%w", err)
	}

	options.Output = "export_sync_remove_from_history_results.json"

	printer.Println("write cleanup result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

	printer.Println("add to history")

	toHistory := m.common.CreateItemsToAdd(items)

	addResult, err := m.syncAddToHistory(client, options, &toHistory)
	if err != nil {
		return fmt.Errorf("add to history error:%w", err)
	}

	options.Output = "export_sync_add_to_history_results.json"

	print("write cleanup result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(addResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (SyncAddToHistoryHandler) syncRemoveFromHistory(client *internal.Client, options *str.Options, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	result, err := client.Sync.RemoveItemsFromHistory(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (SyncAddToHistoryHandler) syncAddToHistory(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.AddResult, error) {
	result, err := client.Sync.AddItemsToHistory(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
