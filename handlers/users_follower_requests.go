// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFollowerRequestsHandler struct for handler
type UsersFollowerRequestsHandler struct{ common CommonLogic }

// Handle to handle users: follower_requests action
func (u UsersFollowerRequestsHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.FollowerRequest == consts.ZeroValue {
		return u.HandleFollowerRequests(options, client)
	}

	if options.Deny {
		return u.HandleFollowerRequestsDeny(options, client)
	}

	return u.HandleApprove(options, client)
}

// HandleApprove approve follower request by id.
func (u UsersFollowerRequestsHandler) HandleApprove(options *str.Options, client *internal.Client) error {
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

// HandleFollowerRequestsDeny deny follower request by id.
func (u UsersFollowerRequestsHandler) HandleFollowerRequestsDeny(options *str.Options, client *internal.Client) error {
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

// HandleFollowerRequests get follower requests.
func (u UsersFollowerRequestsHandler) HandleFollowerRequests(options *str.Options, client *internal.Client) error {
	printer.Println("get follow requests")
	items, err := u.common.FetchFollowRequests(client, options)
	if err != nil {
		return fmt.Errorf("get follow requests error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}
