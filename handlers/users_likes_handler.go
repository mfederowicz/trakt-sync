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

// UsersLikesHandler struct for handler
type UsersLikesHandler struct{ common CommonLogic }

// Handle to handle users: likes action
func (u UsersLikesHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckTypes(options)
	if err != nil {
		return err
	}

	if options.Type != "" {
		printer.Println("Returns all items in a user's likes filtered by type:", options.Type)
	} else {
		printer.Println("Returns all items in a user's likes")
	}

	items, err := u.fetchLikes(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get likes error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersLikesHandler) fetchLikes(client *internal.Client, options *str.Options, page int) ([]*str.UserLike, error) {
	likes, err := u.common.FetchUsersLikes(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch likes error:%w", err)
	}

	if len(likes) == consts.ZeroValue {
		return nil, errors.New("empty likes")
	}

	return likes, nil
}
