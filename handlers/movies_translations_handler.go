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

// MoviesTranslationsHandler struct for handler
type MoviesTranslationsHandler struct{}

// Handle to handle movies: translations action
func (m MoviesTranslationsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all movie translations")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesTranslations(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found translations for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesTranslationsHandler) fetchMoviesTranslations(client *internal.Client, options *str.Options) ([]*str.Translation, *str.Response, error) {
	translations, resp, err := client.Movies.GetAllMovieTranslations(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Language,
	)

	if err != nil {
		return nil, nil, err
	}

	return translations, resp, nil
}
