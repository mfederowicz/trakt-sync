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

// UsersPendingFollowingRequestsHandler struct for handler
type UsersPendingFollowingRequestsHandler struct{ common CommonLogic }

// Handle to handle users: pending_following_requests action
func (u UsersPendingFollowingRequestsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get pending following request")
	items, err := u.common.FetchPendingFollowingRequests(client, options)
	if err != nil {
		return fmt.Errorf("get penging following request error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}
