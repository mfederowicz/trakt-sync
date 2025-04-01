// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// NotesItemHandler struct for handler
type NotesItemHandler struct{ common CommonLogic }

// Handle to handle notes: item action
func (n NotesItemHandler) Handle(options *str.Options, client *internal.Client) error {
if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyNotesIDMsg)
	}

	result, err := n.common.FetchNotesItem(client, options)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}
