// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsRatingsHandler struct for handler
type ShowsRatingsHandler struct{ common CommonLogic }

// Handle to handle shows: ratings action
func (m ShowsRatingsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns rating (between 0 and 10) and distribution for a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsRatings(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found ratings for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsRatingsHandler) fetchShowsRatings(client *internal.Client, options *str.Options) (*str.ShowRatings, *str.Response, error) {
	result, resp, err := client.Shows.GetShowRatings(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
