// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// YounifyRefreshHandler struct for handler
type YounifyRefreshHandler struct{}

// Handle to handle younify: refresh action
func (YounifyRefreshHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validServiceID(options); err != nil {
		return err
	}

	kind := "incremental"
	if options.AllData {
		kind = "full"
	}
	printer.Println("Queue a " + kind + " re-sync of: " + options.ServiceID)
	resp, err := client.Younify.RefreshService(client.BuildCtxFromOptions(options), &options.ServiceID, options.AllData)
	if err = younifyError(options.Action, options.ServiceID, resp, err); err != nil {
		return err
	}

	printer.Printf("result: success, %s re-sync queued for:%s\n", kind, options.ServiceID)
	return nil
}
