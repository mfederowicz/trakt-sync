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

// CommentsItemHandler struct for handler
type CommentsItemHandler struct{ common CommonLogic }

// Handle to handle comments: comments action
func (h CommentsItemHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	result, err := h.common.FetchCommentItem(client, options)
	if err != nil {
		return fmt.Errorf("fetch comment item error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}
