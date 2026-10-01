// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ListsListHandler struct for handler
type ListsListHandler struct{}

// Handle to handle lists: list action
func (h ListsListHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyListIDMsg)
	}

	result, resp, err := h.fetchSingleList(client, options)
	if resp == nil {
		return fmt.Errorf("fetch list error: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found list for:%s", options.InternalID)
	}

	if err != nil {
		return fmt.Errorf("fetch list error: %w", err)
	}

	printer.Printf("Found list for traktId:%s and name:%s \n", options.InternalID, *result.Name)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ListsListHandler) fetchSingleList(client *trakt.Client, options *str.Options) (*str.PersonalList, *str.Response, error) {
	listID := options.InternalID
	result, resp, err := client.Lists.GetList(
		cli.ContextFromOptions(options),
		listID,
	)

	return result, resp, err
}
