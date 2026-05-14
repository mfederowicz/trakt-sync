// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFollowRequestHandler struct for handler
type UsersFollowRequestHandler struct{ common CommonLogic }

// Handle to handle users: follow_request action
func (u UsersFollowRequestHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.FollowerRequest == consts.ZeroValue {
		return errors.New(consts.EmptyFollowerRequestMsg)
	}

	if options.Deny {
		return u.HandleDeny(options, client)
	}

	return u.HandleApprove(options, client)
}

// HandleApprove approve follower request by id.
func (u UsersFollowRequestHandler) HandleApprove(options *str.Options, client *internal.Client) error {
	result, resp, err := u.common.ApproveFollowRequest(client, options)
	if err != nil {
		return fmt.Errorf("approve follower error:%w", err)
	}

	if resp.StatusCode == http.StatusOK {
		printer.Printf("result: success, approve follower:%d \n", options.FollowerRequest)
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		writer.WriteJSON(options, jsonData)
	}

	return nil
}

// HandleDeny deny follower request by id.
func (u UsersFollowRequestHandler) HandleDeny(options *str.Options, client *internal.Client) error {
	result, resp, err := u.common.DenyFollowRequest(client, options)
	if err != nil {
		return fmt.Errorf("deny follower error:%w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, deny follower:%d \n", options.FollowerRequest)
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		writer.WriteJSON(options, jsonData)
	}

	return nil
}
