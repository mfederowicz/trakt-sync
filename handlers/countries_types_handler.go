// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CountriesTypesHandler interface to handle countries types
type CountriesTypesHandler struct{}

// Handle to handle countries: shows action
func (CountriesTypesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("countries handler:" + options.Type)

	countries, _, err := fetchCountries(client, options)
	if err != nil {
		return fmt.Errorf("fetch countries error:%w", err)
	}

	printer.Print("Found " + options.Type + " data \n")
	jsonData, err := json.MarshalIndent(countries, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchCountries(client *trakt.Client, options *str.Options) ([]*str.Country, *str.Response, error) {
	results, resp, err := client.Countries.GetCountries(cli.ContextFromOptions(options), options.Type)

	return results, resp, err
}
