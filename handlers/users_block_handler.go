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

// UsersBlockHandler struct for handler
type UsersBlockHandler struct{ common CommonLogic }

// Handle to handle users: block action
func (m UsersBlockHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.UserName) == consts.ZeroValue {
		return errors.New(consts.EmptyUserNameMsg)
	}

	resp, err := m.usersBlock(client, options)

	if err != nil {
		return fmt.Errorf("block error:%w", err)
	}
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("block success for:", options.UserName)
	}
	return nil
}

func (UsersBlockHandler) usersBlock(client *internal.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Users.Block(client.BuildCtxFromOptions(options), &options.UserName)

	if resp.StatusCode == http.StatusNotFound {
		return resp, fmt.Errorf("user not found:%s", options.UserName)
	}

	if resp.StatusCode == http.StatusConflict {
		return resp, fmt.Errorf("block error:%s", consts.UserBlockedAlready)
	}
	if err != nil {
		return resp, fmt.Errorf("block error:%w", err)
	}

	return resp, nil
}
