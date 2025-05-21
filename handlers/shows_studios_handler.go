// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsStudiosHandler struct for handler
type ShowsStudiosHandler struct{}

// Handle to handle shows: studios action
func (m ShowsStudiosHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all studios for a show")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsStudios(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found studios for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsStudiosHandler) fetchShowsStudios(client *internal.Client, options *str.Options) ([]*str.Studio, *str.Response, error) {
	result, resp, err := client.Shows.GetShowStudios(
		context.Background(),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
