// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// YounifyDisconnectHandler struct for handler
type YounifyDisconnectHandler struct{}

// Handle to handle younify: disconnect action
func (YounifyDisconnectHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validServiceID(options); err != nil {
		return err
	}

	printer.Println("Unlink streaming service: " + options.ServiceID)
	resp, err := client.Younify.DisconnectService(client.BuildCtxFromOptions(options), &options.ServiceID)
	if err = younifyError(options.Action, options.ServiceID, resp, err); err != nil {
		return err
	}

	printer.Println("result: success, unlinked:" + options.ServiceID)
	return nil
}
