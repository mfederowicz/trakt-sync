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
	"github.com/mfederowicz/trakt-sync/writer"
)

// EpisodesTranslationsHandler struct for handler
type EpisodesTranslationsHandler struct{}

// Handle to handle seasons: translations action
func (m EpisodesTranslationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all translations for a specific season of a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s season:%d\n", options.InternalID, options.Episode)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesTranslationsHandler) fetchEpisodesTranslations(client *trakt.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	result, resp, err := client.Shows.GetAllEpisodeTranslations(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Season,
		options.Episode,
		options.Language,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
