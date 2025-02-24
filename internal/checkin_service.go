// Package internal used for client and services
package internal

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinService  handles communication with the checkin related
// methods of the Trakt API.
type CheckinService Service

// DeleteAnyActiveCheckins Removes any active checkins, no need to provide a specific item.
//
// API docs: https://trakt.docs.apiary.io/#reference/checkin/checkin/delete-any-active-checkins
func (c *CheckinService) DeleteAnyActiveCheckins(ctx context.Context) (*str.Response, error) {
	var url = "checkin"
	printer.Println("delete any active checkins")
	req, err := c.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
