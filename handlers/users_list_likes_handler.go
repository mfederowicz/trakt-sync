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

// UsersListLikesHandler struct for handler
type UsersListLikesHandler struct{ common CommonLogic }

// Handle to handle users: list_likes action
func (u UsersListLikesHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Returns all users who liked a list.")
	items, err := u.fetchListLikes(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get list likes error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersListLikesHandler) fetchListLikes(client *internal.Client, options *str.Options, page int) ([]*str.UserLike, error) {
	likes, err := u.common.FetchUsersListLikes(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch likes error:%w", err)
	}

	if len(likes) == consts.ZeroValue {
		return nil, errors.New("empty likes")
	}

	return likes, nil
}
