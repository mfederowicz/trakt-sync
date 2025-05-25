// Package handlers used to handle module actions
package handlers

import (
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

// MoviesBoxofficeHandler struct for handler
type MoviesBoxofficeHandler struct{}

// Handle to handle movies: boxoffice action
func (h MoviesBoxofficeHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns the top 10 grossing movies in the U.S. box office last weekend. Updated every Monday morning.")
	result, err := h.fetchMoviesBoxoffice(client, options)
	if err != nil {
		return fmt.Errorf("fetch movies error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty movies")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.MoviesItem{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (MoviesBoxofficeHandler) fetchMoviesBoxoffice(client *internal.Client, options *str.Options) ([]*str.MoviesItem, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Movies.GetBoxoffice(
		client.BuildCtxFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}
