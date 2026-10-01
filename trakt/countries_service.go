package trakt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
)

// CountriesService  handles communication with the countries related
// methods of the Trakt API.
type CountriesService Service

// GetCountries Get a list of all countries, including names and codes.
//
// API docs: https://docs.trakt.tv/reference/getcountrieslist
func (c *CountriesService) GetCountries(ctx context.Context, strType string) ([]*str.Country, *str.Response, error) {
	var url = fmt.Sprintf("countries/%s", strType)
	c.client.debug("fetch countries url:" + url)

	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	list := []*str.Country{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		c.client.debug("fetch countries err:", err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
