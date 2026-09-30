// Package handlers used to handle module actions
package handlers

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CheckinDeleteHandler struct for handler
type CheckinDeleteHandler struct{}

// Handle to handle checkin: episode action
func (h CheckinDeleteHandler) Handle(options *str.Options, client *trakt.Client) error {
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
func (CheckinDeleteHandler) deleteActiveCheckins(client *trakt.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Checkin.DeleteAnyActiveCheckins(
		cli.ContextFromOptions(options),
	)

	return resp, err
}
