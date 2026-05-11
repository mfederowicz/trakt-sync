// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncGetHistoryHandler struct for handler
type SyncGetHistoryHandler struct{ common CommonLogic }

// Handle to handle sync: get_history action
func (m SyncGetHistoryHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	err = m.common.CheckDates(options.StartDate, options.EndDate, options.Timezone)
	if err != nil {
		return err
	}

	printer.Println("Get watched history type:", options.Type)
	items, err := m.syncGetHistoryItems(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get watched error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m SyncGetHistoryHandler) syncGetHistoryItems(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := m.common.FetchHistoryList(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
