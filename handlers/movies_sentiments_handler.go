// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesSentimentsHandler struct for handler
type MoviesSentimentsHandler struct{}

// Handle to handle movies: sentiments action
func (MoviesSentimentsHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	printer.Println("Returns sentiment counts for comments and reactions attached to a movie.")
	result, resp, err := client.Movies.GetMovieSentiments(client.BuildCtxFromOptions(options), &options.InternalID)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found movie for:%s", options.InternalID)
	}
	if err != nil {
		return fmt.Errorf("fetch sentiments error: %w", err)
	}

	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal sentiments error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
