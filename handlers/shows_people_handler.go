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

// ShowsPeopleHandler struct for handler
type ShowsPeopleHandler struct{ common CommonLogic }

// Handle to handle shows: people action
func (m ShowsPeopleHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all people related with show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowPeople(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found people for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

// fetchShowPeople to fetch all people from show
func (ShowsPeopleHandler) fetchShowPeople(client *internal.Client, options *str.Options) (*str.ShowPeople, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.People.GetAllPeopleForShow(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
