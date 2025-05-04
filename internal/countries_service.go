// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CountriesService  handles communication with the countries related
// methods of the Trakt API.
type CountriesService Service

// GetCountries Get a list of all countries, including names and codes.
//
// API docs: https://trakt.docs.apiary.io/#reference/countries/list/get-countries
func (c *CountriesService) GetCountries(ctx context.Context, strType *string) ([]*str.Country, *str.Response, error) {
	var url = fmt.Sprintf("countries/%s", *strType)
	printer.Println("fetch countries url:" + url)

	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	list := []*str.Country{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch countries err:", err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
