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

// CommentsCommentHandler struct for handler
type CommentsCommentHandler struct{ common CommonLogic }

// Handle to handle comments: comment action
func (h CommentsCommentHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	if options.Delete {
		return h.HandleDelete(options, client)
	}

	if len(options.Comment) > consts.ZeroValue {
		return h.HandleModify(options, client)
	}

	result, err := h.common.FetchComment(client, options)
	if err != nil {
		return fmt.Errorf("fetch comment error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

// HandleModify modify exiting comment
func (h CommentsCommentHandler) HandleModify(options *str.Options, client *internal.Client) error {
	c := new(str.Comment)
	c.Comment = &options.Comment
	c.Spoiler = &options.Spoiler
	result, resp, err := h.common.UpdateComment(client, options, c)
	if err != nil {
		return fmt.Errorf("update comment error:%w", err)
	}

	if resp.StatusCode == http.StatusOK {
		printer.Printf("result: success, update comment:%d \n", result.ID)
	}

	return nil
}

// HandleDelete handle delete comment
func (h CommentsCommentHandler) HandleDelete(options *str.Options, client *internal.Client) error {
	resp, err := h.common.DeleteComment(client, options)
	if err != nil {
		return fmt.Errorf("delete comment error:%w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, remove comment:%d \n", options.CommentID)
	}
	return nil
}
