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

// CheckintoAnItem Check into a movie or episode.
//
// API docs: https://trakt.docs.apiary.io/#reference/checkin/checkin/check-into-an-item
func (c *CheckinService) CheckintoAnItem(ctx context.Context, checkin *str.CheckIn) (*str.CheckIn, *str.Response, error) {
	var url = "checkin"
	printer.Println("create new checkin")
	req, err := c.client.NewRequest(http.MethodPost, url, checkin)
	if err != nil {
		return nil, nil, err
	}

	ch := new(str.CheckIn)
	resp, err := c.client.Do(ctx, req, ch)
	if err != nil {
		return ch, resp, err
	}

	return ch, resp, nil
}
