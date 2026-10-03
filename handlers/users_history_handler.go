// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersHistoryHandler struct for handler
type UsersHistoryHandler struct{ common CommonLogic }

// Handle to handle users: history action
func (m UsersHistoryHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := checkMediaFilters(options); err != nil {
		return err
	}
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}
	if options.ItemID > consts.ZeroValue && len(options.Type) == consts.ZeroValue {
		return errors.New(consts.EmptyHistoryItemTypeMsg)
	}

	err = m.common.CheckDates(options.StartDate, options.EndDate, options.Timezone)
	if err != nil {
		return err
	}

	printer.Println("Get watched history type:", options.Type)
	items, err := m.usersHistory(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get watched error:%w", err)
	}
	if len(items) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersHistoryHandler) usersHistory(client *trakt.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := m.common.FetchUsersHistory(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
