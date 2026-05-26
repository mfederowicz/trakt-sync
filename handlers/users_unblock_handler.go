// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersUnblockHandler struct for handler
type UsersUnblockHandler struct{ common CommonLogic }

// Handle to handle users: unblock action
func (m UsersUnblockHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.UserName) == consts.ZeroValue {
		return errors.New(consts.EmptyUserNameMsg)
	}

	resp, err := m.usersUnblock(client, options)

	if err != nil {
		return fmt.Errorf("unblock error:%w", err)
	}
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("unblock success for:", options.UserName)
	}

	return nil
}

func (UsersUnblockHandler) usersUnblock(client *internal.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Users.Unblock(client.BuildCtxFromOptions(options), &options.UserName)

	if resp.StatusCode == http.StatusNotFound {
		return resp, fmt.Errorf("user not found:%s", options.UserName)
	}

	if err != nil {
		return resp, fmt.Errorf("unblock error:%w", err)
	}

	return resp, nil
}
