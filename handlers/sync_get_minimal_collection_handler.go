// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncGetMinimalCollectionHandler struct for handler
type SyncGetMinimalCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: get_minimal_collection action
func (s SyncGetMinimalCollectionHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := s.common.CheckTypes(options); err != nil {
		return err
	}
	if len(options.AvailableOn) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.SyncAvailableOn, options.AvailableOn) {
		return fmt.Errorf("available_on '%s' is not valid, avaliable values: %v", options.AvailableOn, cfg.SyncAvailableOn)
	}

	printer.Println("Get minimal collection type:", options.Type)
	result, err := s.fetchMinimalCollection(client, options)
	if err != nil {
		return fmt.Errorf("get minimal collection error: %w", err)
	}

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode minimal collection: %w", err)
	}

	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncGetMinimalCollectionHandler) fetchMinimalCollection(client *trakt.Client, options *str.Options) (any, error) {
	ctx := cli.ContextFromOptions(options)
	opts := uri.ListOptions{AvailableOn: options.AvailableOn}
	if options.Type == consts.Shows {
		result, _, err := client.Sync.GetMinimalShowCollection(ctx, &opts)
		return result, err
	}

	result, _, err := client.Sync.GetMinimalCollection(ctx, &options.Type, &opts)
	return result, err
}
