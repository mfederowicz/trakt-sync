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
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsSummaryHandler struct for handler
type ShowsSummaryHandler struct{}

// Handle to handle shows: summary action
func (m ShowsSummaryHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns a single show details")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsSummary(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found show for id:%s and name:%s \n", options.InternalID, *result.Title)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsSummaryHandler) fetchShowsSummary(client *internal.Client, options *str.Options) (*str.Show, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	show, resp, err := client.Shows.GetShow(
		context.Background(),
		&options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return show, resp, nil
}
