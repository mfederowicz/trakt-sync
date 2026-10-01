// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersBlockedUsersHandler struct for handler
type UsersBlockedUsersHandler struct{ common CommonLogic }

// Handle to handle users: blocked_users action
func (u UsersBlockedUsersHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all users you have blocked, including when each user was blocked.")

	items, err := u.fetchBlockedUsers(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get blocked users error:%w", err)
	}
	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersBlockedUsersHandler) fetchBlockedUsers(client *trakt.Client, options *str.Options, page int) ([]*str.UserBlocked, error) {
	users, err := u.common.FetchBlockedUsers(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch blocked users error:%w", err)
	}

	if len(users) == consts.ZeroValue {
		return nil, errors.New("empty list")
	}

	return users, nil
}
