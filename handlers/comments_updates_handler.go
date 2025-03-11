// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CommentsUpdatesHandler struct for handler
type CommentsUpdatesHandler struct{ common CommonLogic }

// Handle to handle comments: updates action
func (h CommentsUpdatesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns the most recently updated comments across all of Trakt.")

	result, err := h.common.FetchUpdatedComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch comments error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty list")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.CommentItem{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}
