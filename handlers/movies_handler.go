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

// MoviesHandler interface to handle movies module action
type MoviesHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// writeMoviesItems writes a fetched movies list to the output file.
func writeMoviesItems(options *str.Options, result []*str.MoviesItem) error {
	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal movies %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
