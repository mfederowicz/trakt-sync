// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// YounifyConnectHandler struct for handler
type YounifyConnectHandler struct{}

// Handle to handle younify: connect action
func (YounifyConnectHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validServiceID(options); err != nil {
		return err
	}

	returnURL := options.ReturnURL
	if len(returnURL) == consts.ZeroValue {
		returnURL = consts.DefaultReturnURL
	}

	printer.Println("Create a streaming connection for: " + options.ServiceID)
	connect := &str.YounifyConnect{ServiceID: &options.ServiceID, ReturnURL: &returnURL}
	result, resp, err := client.Younify.Connect(client.BuildCtxFromOptions(options), connect)
	if err = younifyError(options.Action, options.ServiceID, resp, err); err != nil {
		return err
	}

	if result.URL == nil || len(*result.URL) == consts.ZeroValue {
		return errors.New("younify connect: no web auth URL in the response")
	}

	printer.Println("Open this URL to connect " + options.ServiceID + ": " + *result.URL)
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal younify connect error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
