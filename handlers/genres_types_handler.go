// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// GenresTypesHandler interface to handle genres types
type GenresTypesHandler struct{}

// Handle to handle genres: shows action
func (GenresTypesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("genres handler:" + options.Type)

	genres, _, err := fetchGenres(client, &options.Type)
	if err != nil {
		return fmt.Errorf("fetch genres error:%w", err)
	}

	printer.Print("Found " + options.Type + " data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(genres, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchGenres(client *internal.Client, strType *string) ([]*str.Genre, *str.Response, error) {
	results, resp, err := client.Genres.GetGenres(context.Background(), strType)

	return results, resp, err
}
