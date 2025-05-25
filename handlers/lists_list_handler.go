// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ListsListHandler struct for handler
type ListsListHandler struct{}

// Handle to handle lists: list action
func (h ListsListHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyListIDMsg)
	}

	result, resp, _ := h.fetchSingleList(client, options)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found list for:%s", options.InternalID)
	}

	printer.Printf("Found list for traktId:%s and name:%s \n", options.InternalID, *result.Name)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ListsListHandler) fetchSingleList(client *internal.Client, options *str.Options) (*str.PersonalList, *str.Response, error) {
	listID := options.InternalID
	result, resp, err := client.Lists.GetList(
		client.BuildCtxFromOptions(options),
		&listID,
	)

	return result, resp, err
}
