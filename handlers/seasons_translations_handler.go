// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SeasonsTranslationsHandler struct for handler
type SeasonsTranslationsHandler struct{}

// Handle to handle seasons: translations action
func (m SeasonsTranslationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all translations for a specific season of a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s season:%d\n", options.InternalID, options.Season)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsTranslationsHandler) fetchSeasonsTranslations(client *trakt.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetAllSeasonTranslations(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Season,
		options.Language,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
