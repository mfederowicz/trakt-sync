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

// ShowsStatsHandler struct for handler
type ShowsStatsHandler struct{ common CommonLogic }

// Handle to handle shows: stats action
func (m ShowsStatsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns lots of show stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found stats for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsStatsHandler) fetchShowsStats(client *internal.Client, options *str.Options) (*str.ShowStats, *str.Response, error) {
	result, resp, err := client.Shows.GetShowStats(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
