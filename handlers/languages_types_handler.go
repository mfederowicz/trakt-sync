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

// LanguagesTypesHandler interface to handle languages types
type LanguagesTypesHandler struct{}

// Handle to handle languages: shows action
func (LanguagesTypesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("languages handler:" + options.Type)

	languages, _, err := fetchLanguages(client, &options.Type)
	if err != nil {
		return fmt.Errorf("fetch languages error:%w", err)
	}

	printer.Print("Found " + options.Type + " data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(languages, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchLanguages(client *internal.Client, strType *string) ([]*str.Language, *str.Response, error) {
	results, resp, err := client.Languages.GetLanguages(context.Background(), strType)

	return results, resp, err
}
