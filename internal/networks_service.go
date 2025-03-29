// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// NetworksService  handles communication with the movies related
// methods of the Trakt API.
type NetworksService Service

// GetNetworksList Get a list of all TV networks, including the name, country, and ids.
func (m *NetworksService) GetNetworksList(ctx context.Context, opts *uri.ListOptions) ([]*str.TvNetwork, *str.Response, error) {
	var url = "networks"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch networks url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.TvNetwork{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, errors.New("not found networks")
	}

	if err != nil {
		printer.Println("fetch networks err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
