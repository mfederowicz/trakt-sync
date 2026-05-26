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

// UsersUnfollowHandler struct for handler
type UsersUnfollowHandler struct{ common CommonLogic }

// Handle to handle users: unfollow action
func (m UsersUnfollowHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.UserName) == consts.ZeroValue {
		return errors.New(consts.EmptyUserNameMsg)
	}

	resp, err := m.usersUnfollow(client, options)

	if err != nil {
		return fmt.Errorf("unfollow error:%w", err)
	}
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("unfollow success for:", options.UserName)
	}

	return nil
}

func (UsersUnfollowHandler) usersUnfollow(client *internal.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Users.Unfollow(client.BuildCtxFromOptions(options), &options.UserName)

	if resp.StatusCode == http.StatusNotFound {
		return resp, fmt.Errorf("user not found:%s", options.UserName)
	}

	if err != nil {
		return resp, fmt.Errorf("unfollow error:%w", err)
	}

	return resp, nil
}
