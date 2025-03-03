// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CommentsRepliesHandler struct for handler
type CommentsRepliesHandler struct{ common CommonLogic }

// Handle to handle comments: replies action
func (h CommentsRepliesHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	if len(options.Reply) > consts.ZeroValue {
		return h.replyForComment(client, options)
	}

	return h.allCommentReplies(client, options)
}

func (h CommentsRepliesHandler) allCommentReplies(client *internal.Client, options *str.Options) error {
	printer.Println("Returns all replies for a comment.")
	result, err := h.fetchCommentReplies(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch replies error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty replies")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.Comment{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func (h CommentsRepliesHandler) fetchCommentReplies(client *internal.Client, options *str.Options, page int) ([]*str.Comment, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Comments.GetRepliesForComment(
		context.Background(),
		&opts,
		&options.CommentID,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := h.fetchCommentReplies(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

func (h CommentsRepliesHandler) replyForComment(client *internal.Client, options *str.Options) error {
	c := new(str.Comment)
	c.Comment = &options.Reply
	c.Spoiler = &options.Spoiler
	result, resp, err := h.common.Reply(client, &options.CommentID, c)
	if err != nil {
		return fmt.Errorf("reply comment error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, reply comment:%d \n", result.ID)
	}

	return nil
}
