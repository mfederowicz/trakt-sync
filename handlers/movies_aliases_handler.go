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

// MoviesAliasesHandler struct for handler
type MoviesAliasesHandler struct{}

// Handle to handle people: aliases action
func (m MoviesAliasesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns a single movie details")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesAliases(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found aliases for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesAliasesHandler) fetchMoviesAliases(client *trakt.Client, options *str.Options) ([]*str.Alias, *str.Response, error) {
	aliases, resp, err := client.Movies.GetAllMovieAliases(
		cli.ContextFromOptions(options),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return aliases, resp, nil
}
