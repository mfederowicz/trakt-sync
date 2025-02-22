// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersSavedFiltersHandler struct for handler
type UsersSavedFiltersHandler struct{}

// Handle to handle users: saved filters action
func (UsersSavedFiltersHandler) Handle(options *str.Options, client *internal.Client) error {	
	printer.Println("users saved filters handler:" + options.UserName)

	filters, resp, err := fetchUsersSavedFilters(client, &options.Type)
	if err != nil {
		return fmt.Errorf("fetch saved filters error:%w", err)
	}

	if resp.StatusCode == http.StatusUpgradeRequired {
		return cli.HandleUpgrade(resp)	
	}

	if len(filters) == consts.ZeroValue {
		return errors.New("empty list of filters")
	}

	printer.Print("Found " + options.Action + " data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(filters, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchUsersSavedFilters(client *internal.Client, section *string) ([]*str.SavedFilter, *str.Response, error) {
	lists, resp, err := client.Users.GetSavedFilters(
		context.Background(),
		section,
	)

	return lists, resp, err
}

