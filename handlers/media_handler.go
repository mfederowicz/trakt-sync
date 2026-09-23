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

// MediaHandler interface to handle media module action
type MediaHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// exportMedia fetches all pages of a media list and writes them to the output file.
func exportMedia[T any](client *internal.Client, options *str.Options, fetch pageFetcher[T]) error {
	result, err := fetchAllPages(client, options, consts.DefaultPage, fetch)
	if err != nil {
		return fmt.Errorf("fetch media %s error: %w", options.Action, err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal media %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
