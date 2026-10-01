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

// MoviesTranslationsHandler struct for handler
type MoviesTranslationsHandler struct{}

// Handle to handle movies: translations action
func (m MoviesTranslationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all movie translations")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesTranslationsHandler) fetchMoviesTranslations(client *trakt.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	translations, resp, err := client.Movies.GetAllMovieTranslations(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Language,
	)

	if err != nil {
		return nil, nil, err
	}

	return translations, resp, nil
}
