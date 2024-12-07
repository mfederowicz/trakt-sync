// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// PeopleSummaryHandler struct for handler
type PeopleSummaryHandler struct{}

// Handle to handle people: summary action
func (p PeopleSummaryHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New("set personId ie: -i john-wayne")
	}
	printer.Println("Get a single person")
	result, err := p.fetchSinglePerson(client, options)
	if err != nil {
		return fmt.Errorf("fetch single person error:%w", err)
	}

	if result == nil {
		return errors.New("empty result")
	}

	printer.Print("Found person \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (PeopleSummaryHandler) fetchSinglePerson(client *internal.Client, options *str.Options) (*str.Person, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.People.GetSinglePerson(
		context.Background(),
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
