// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsTranslationsHandler struct for handler
type ShowsTranslationsHandler struct{}

// Handle to handle shows: translations action
func (m ShowsTranslationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all show translations")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsTranslationsHandler) fetchShowsTranslations(client *trakt.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	translations, resp, err := client.Shows.GetAllShowTranslations(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Language,
	)

	if err != nil {
		return nil, nil, err
	}

	return translations, resp, nil
}
