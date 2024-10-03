// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// PeopleShowsHandler struct for handler
type PeopleShowsHandler struct{}

// Handle to handle people: shows action
func (p PeopleShowsHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return fmt.Errorf(consts.EmptyPersonIDMsg)
	}
	printer.Println("Get show credits")
	result, err := p.fetchShowCredits(client, options)
	if err != nil {
		return fmt.Errorf("fetch show credits error:%w", err)
	}

	if result == nil {
		return fmt.Errorf(consts.EmptyResult)
	}

	printer.Print("Found show credits data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

 func (PeopleShowsHandler) fetchShowCredits(client *internal.Client, options *str.Options) (*str.PersonShows, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.People.GetShowCredits(
		context.Background(),
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

