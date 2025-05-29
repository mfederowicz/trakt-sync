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

// SeasonsTranslationsHandler struct for handler
type SeasonsTranslationsHandler struct{}

// Handle to handle seasons: translations action
func (m SeasonsTranslationsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all translations for a specific season of a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s season:%d\n", options.InternalID, options.Season)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsTranslationsHandler) fetchSeasonsTranslations(client *internal.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetAllSeasonTranslations(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
		&options.Language,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
