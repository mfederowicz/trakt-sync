// Package handlers used to handle module actions
package handlers

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinDeleteHandler struct for handler
type CheckinDeleteHandler struct{}

// Handle to handle checkin: episode action
func (h CheckinDeleteHandler) Handle(options *str.Options, client *internal.Client) error {
	resp, err := h.deleteActiveCheckins(client, options)
	if resp == nil {
		return fmt.Errorf("delete checkin error: %w", err)
	}

	if err != nil {
		return fmt.Errorf("delete checkin error: %w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Print("result: success \n")
	}

	return nil
}
func (CheckinDeleteHandler) deleteActiveCheckins(client *internal.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Checkin.DeleteAnyActiveCheckins(
		client.BuildCtxFromOptions(options),
	)

	return resp, err
}
