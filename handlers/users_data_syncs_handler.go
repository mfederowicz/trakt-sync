// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// UsersDataSyncsHandler struct for handler
type UsersDataSyncsHandler struct{}

// Handle to handle users: data_syncs action
func (UsersDataSyncsHandler) Handle(options *str.Options, client *internal.Client) error {
	if !cfg.IsValidConfigType(cfg.DataSyncTypes, options.Type) {
		return fmt.Errorf("set -t to one of %v, or leave it out for all data syncs", cfg.DataSyncTypes)
	}

	printer.Println("Returns your data syncs " + options.Type)
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.DataSync, *str.Response, error) {
		return client.Users.GetDataSyncs(client.BuildCtxFromOptions(options), &options.Type, opts)
	})
	if apiErr := notOpenToAPIApps(options.Action, err); apiErr != nil {
		return apiErr
	}
	if err != nil {
		return fmt.Errorf("fetch data syncs error: %w", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeResult(options, result)
}

// UsersDataSyncHandler struct for handler
type UsersDataSyncHandler struct{}

// Handle to handle users: data_sync action
func (UsersDataSyncHandler) Handle(options *str.Options, client *internal.Client) error {
	id, err := dataSyncID(options)
	if err != nil {
		return err
	}

	printer.Printf("Returns data sync %d\n", id)
	result, resp, err := client.Users.GetDataSync(client.BuildCtxFromOptions(options), id)
	if err = dataSyncError(options.Action, id, resp, err); err != nil {
		return err
	}

	return writeResult(options, result)
}

// UsersDataSyncItemsHandler struct for handler
type UsersDataSyncItemsHandler struct{}

// Handle to handle users: data_sync_paused and data_sync_skipped actions
func (UsersDataSyncItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	id, err := dataSyncID(options)
	if err != nil {
		return err
	}

	section := strings.TrimPrefix(options.Action, consts.DataSync+"_")
	printer.Printf("Returns %s items of data sync %d\n", section, id)
	var lastResp *str.Response
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.SyncItem, *str.Response, error) {
		list, resp, fetchErr := client.Users.GetDataSyncItems(client.BuildCtxFromOptions(options), id, &section, opts)
		lastResp = resp
		return list, resp, fetchErr
	})
	if err = dataSyncError(options.Action, id, lastResp, err); err != nil {
		return err
	}

	if len(result) == consts.ZeroValue {
		return fmt.Errorf("no %s items in data sync %d", section, id)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeResult(options, result)
}

// UsersUndoDataSyncHandler struct for handler
type UsersUndoDataSyncHandler struct{}

// Handle to handle users: undo_data_sync action
func (UsersUndoDataSyncHandler) Handle(options *str.Options, client *internal.Client) error {
	id, err := dataSyncID(options)
	if err != nil {
		return err
	}

	printer.Printf("Undo data sync %d: reverses every item it imported\n", id)
	resp, err := client.Users.UndoDataSync(client.BuildCtxFromOptions(options), id)
	if err = dataSyncError(options.Action, id, resp, err); err != nil {
		return err
	}

	printer.Printf("result: success, undone data sync:%d\n", id)
	return nil
}

// dataSyncID reads the numeric data sync id from -i.
func dataSyncID(options *str.Options) (int, error) {
	id, err := strconv.Atoi(options.ID)
	if err != nil || id <= consts.ZeroValue {
		return consts.ZeroValue, errors.New("set data sync id ie: -i 157 (ids from users -a data_syncs)")
	}
	return id, nil
}

// dataSyncError maps a data sync response to a readable error; a sync of another user answers 404.
func dataSyncError(action string, id int, resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found data sync:%d", id)
	}
	if apiErr := notOpenToAPIApps(action, err); apiErr != nil {
		return apiErr
	}
	if err != nil {
		return fmt.Errorf("%s error: %w", action, err)
	}
	return nil
}
