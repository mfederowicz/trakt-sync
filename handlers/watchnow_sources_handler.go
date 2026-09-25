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

// WatchNowSourcesHandler struct for handler
type WatchNowSourcesHandler struct{}

// Handle to handle watchnow: sources action
func (WatchNowSourcesHandler) Handle(options *str.Options, client *internal.Client) error {
	ctx := client.BuildCtxFromOptions(options)
	var (
		result []map[string][]*str.WatchNowSource
		resp   *str.Response
		err    error
	)
	country := consts.ActionTypeAll
	if len(options.Country) > consts.ZeroValue {
		country = options.Country
		printer.Println("Returns watch now sources available in: " + country + " (limited access).")
		result, resp, err = client.WatchNow.GetWatchNowSourcesByCountry(ctx, &options.Country)
	} else {
		printer.Println("Returns all watch now sources supported by Trakt (limited access).")
		result, resp, err = client.WatchNow.GetWatchNowSources(ctx)
	}
	if err = watchNowError(consts.WatchNow+" "+consts.Sources, consts.Sources, country, resp, err); err != nil {
		return err
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal watchnow sources error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
