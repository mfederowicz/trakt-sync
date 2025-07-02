// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncAddToCollectionHandler struct for handler
type SyncAddToCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_collection action
func (m SyncAddToCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.ReadInput(options.CollectionItems)
	if err != nil {
		return err
	}
	printer.Println("Add collection")
	result, err := m.syncAddToCollection(client, options, items)
	if err != nil {
		return fmt.Errorf("add to collection error:%w", err)
	}

	print("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

// ReadInput reads data from file or from stdin
func (m SyncAddToCollectionHandler) ReadInput(filePath string) (*str.CollectionItems, error) {
	if filePath != consts.EmptyString {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		return m.common.ConvertBytesToColletionItems(data)
	}

	// Check if there's data in stdin to avoid blocking
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}

	// os.ModeCharDevice means no data is being piped (stdin is a terminal)
	if fi.Mode()&os.ModeCharDevice != 0 {
		return nil, fmt.Errorf("no --file provided and no data piped to stdin")
	}

	// Read all data from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}

	return m.common.ConvertBytesToColletionItems(data)
}

func (SyncAddToCollectionHandler) syncAddToCollection(client *internal.Client, options *str.Options, items *str.CollectionItems) (*str.CollectionAddResult, error) {
	result, err := client.Sync.AddItemsToCollection(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
