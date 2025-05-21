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

// ShowsAliasesHandler struct for handler
type ShowsAliasesHandler struct{}

// Handle to handle shows: aliases action
func (m ShowsAliasesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns a single show details")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsAliases(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found aliases for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsAliasesHandler) fetchShowsAliases(client *internal.Client, options *str.Options) ([]*str.Alias, *str.Response, error) {
	aliases, resp, err := client.Shows.GetAllShowAliases(
		context.Background(),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return aliases, resp, nil
}
