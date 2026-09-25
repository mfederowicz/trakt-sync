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

// YounifyConnectionsHandler struct for handler
type YounifyConnectionsHandler struct{}

// Handle to handle younify: connections action
func (YounifyConnectionsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns streaming services with your connection status.")
	result, resp, err := client.Younify.GetConnections(client.BuildCtxFromOptions(options))
	if err = younifyError(options.Action, consts.EmptyString, resp, err); err != nil {
		return err
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal younify connections error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
