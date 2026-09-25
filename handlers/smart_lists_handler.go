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

// SmartListsHandler interface to handle smart_lists module action
type SmartListsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// smartListError maps a smart list response to a readable error; private lists answer 404 unless you own them.
func smartListError(action string, id string, resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found smart list for:%s (private lists are visible only to their owner)", id)
	}

	if err != nil {
		return fmt.Errorf("fetch smart list %s error: %w", action, err)
	}

	return nil
}

// writeSmartList writes a smart list result to the output file.
func writeSmartList(options *str.Options, data any) error {
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(data, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal smart list %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}

// validSmartListID checks the smart list slug needed by every smart_lists action.
func validSmartListID(options *str.Options) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySmartListIDMsg)
	}
	return nil
}
