// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SeasonsPeopleHandler struct for handler
type SeasonsPeopleHandler struct{ common CommonLogic }

// Handle to handle seasons: people action
func (m SeasonsPeopleHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all people related with season.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsPeople(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found people for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsPeopleHandler) fetchSeasonsPeople(client *internal.Client, options *str.Options) (*str.SeasonPeople, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetAllPeopleForSeason(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
