// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFollowRequestsHandler struct for handler
type UsersFollowRequestsHandler struct{ common CommonLogic }

// Handle to handle users: follow_requests action
func (u UsersFollowRequestsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get follow requests")
	items, err := u.common.FetchFollowRequests(client, options)
	if err != nil {
		return fmt.Errorf("get follow requests error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}
