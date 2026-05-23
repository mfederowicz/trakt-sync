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

// UsersListCommentsHandler struct for handler
type UsersListCommentsHandler struct{ common CommonLogic }

// Handle to handle users: list_comments action
func (u UsersListCommentsHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckSortAndTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Returns all top level comments for a list")

	items, err := u.fetchComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get comments error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersListCommentsHandler) fetchComments(client *internal.Client, options *str.Options, page int) ([]*str.ListComment, error) {
	comments, err := u.common.FetchUsersListComments(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch list comments error:%w", err)
	}

	if len(comments) == consts.ZeroValue {
		return nil, errors.New("empty comments")
	}

	return comments, nil
}
