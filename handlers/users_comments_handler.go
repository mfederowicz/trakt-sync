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

// UsersCommentsHandler struct for handler
type UsersCommentsHandler struct{ common CommonLogic }

// Handle to handle users: comments action
func (u UsersCommentsHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckCommentsFilters(options)
	if err != nil {
		return err
	}

	if options.Type != "" {
		printer.Println("Returns all items in a user's comments filtered by type:", options.Type)
	} else {
		printer.Println("Returns all items in a user's comments")
	}

	items, err := u.fetchComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get likes error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersCommentsHandler) fetchComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error) {
	comments, err := u.common.FetchUsersComments(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch comments error:%w", err)
	}

	if len(comments) == consts.ZeroValue {
		return nil, errors.New("empty comments")
	}

	return comments, nil
}
