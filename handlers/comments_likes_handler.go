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

// CommentsLikesHandler struct for handler
type CommentsLikesHandler struct{ common CommonLogic }

// Handle to handle comments: likes action
func (h CommentsLikesHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	printer.Println("Get all users who liked a comment.")
	result, err := h.common.FetchCommentUserLikes(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch likes error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty list")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.CommentUserLike{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}
