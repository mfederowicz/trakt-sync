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

// EpisodesTranslationsHandler struct for handler
type EpisodesTranslationsHandler struct{}

// Handle to handle seasons: translations action
func (m EpisodesTranslationsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all translations for a specific season of a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s season:%d\n", options.InternalID, options.Episode)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesTranslationsHandler) fetchEpisodesTranslations(client *internal.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	result, resp, err := client.Shows.GetAllEpisodeTranslations(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
		&options.Episode,
		&options.Language,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
