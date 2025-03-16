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

// CountriesTypesHandler interface to handle countries types
type CountriesTypesHandler struct{}

// Handle to handle countries: shows action
func (CountriesTypesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("countries handler:" + options.Type)

	countries, _, err := fetchCountries(client, &options.Type)
	if err != nil {
		return fmt.Errorf("fetch countries error:%w", err)
	}

	printer.Print("Found " + options.Type + " data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(countries, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchCountries(client *internal.Client, strType *string) ([]*str.Country, *str.Response, error) {
	results, resp, err := client.Countries.GetCountries(context.Background(), strType)

	return results, resp, err
}
