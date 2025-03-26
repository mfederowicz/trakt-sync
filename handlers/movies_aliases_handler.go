// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesAliasesHandler struct for handler
type MoviesAliasesHandler struct{}

// Handle to handle people: aliases action
func (m MoviesAliasesHandler) Handle(options *str.Options, client *internal.Client) error {
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

func (MoviesAliasesHandler) fetchMoviesAliases(client *internal.Client, options *str.Options) ([]*str.Alias, *str.Response, error) {
	aliases, resp, err := client.Movies.GetAllMovieAliases(
		context.Background(),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return aliases, resp, nil
}

