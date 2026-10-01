// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFollowHandler struct for handler
type UsersFollowHandler struct{}

// Handle to handle users: follow action
func (m UsersFollowHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.UserName) == consts.ZeroValue {
		return errors.New(consts.EmptyUserNameMsg)
	}

	result, resp, err := m.usersFollow(client, options)

	if err != nil {
		return fmt.Errorf("follow error:%w", err)
	}
	if resp.StatusCode == http.StatusCreated {
		printer.Println("follow success for:", options.UserName)
	}
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (UsersFollowHandler) usersFollow(client *trakt.Client, options *str.Options) (*str.FollowResult, *str.Response, error) {
	result, resp, err := client.Users.Follow(cli.ContextFromOptions(options), options.UserName)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, resp, fmt.Errorf("user not found:%s", options.UserName)
	}

	if resp != nil && resp.StatusCode == http.StatusConflict {
		return nil, resp, fmt.Errorf("follow error:%s", consts.UserPendingFollowRequest)
	}
	if err != nil {
		return nil, resp, fmt.Errorf("follow error:%w", err)
	}

	return result, resp, nil
}
