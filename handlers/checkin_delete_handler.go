// Package handlers used to handle module actions
package handlers

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinDeleteHandler struct for handler
type CheckinDeleteHandler struct{}

// Handle to handle checkin: episode action
func (h CheckinDeleteHandler) Handle(options *str.Options, client *internal.Client) error {
	resp, _ := h.deleteActiveCheckins(client, options)	

	if resp.StatusCode == http.StatusNoContent {
		printer.Print("result: success \n")
	}

	return nil
}
func (CheckinDeleteHandler) deleteActiveCheckins(client *internal.Client, _ *str.Options) (*str.Response, error) {
	resp, err := client.Checkin.DeleteAnyActiveCheckins(
		context.Background(),
	)

	return resp, err
}
