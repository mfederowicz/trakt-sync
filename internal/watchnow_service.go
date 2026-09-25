// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
)

// WatchNowService handles communication with the watch now related
// methods of the Trakt API.
type WatchNowService Service

// GetWatchNowSources Returns all watch now sources supported by Trakt, grouped by country.
//
// API docs: https://docs.trakt.tv/reference/getwatchnowsourcesall
func (w *WatchNowService) GetWatchNowSources(ctx context.Context) ([]map[string][]*str.WatchNowSource, *str.Response, error) {
	return w.fetchSources(ctx, "watchnow/sources")
}

// GetWatchNowSourcesByCountry Returns watch now sources available in a country.
//
// API docs: https://docs.trakt.tv/reference/getwatchnowsourcescountry
func (w *WatchNowService) GetWatchNowSourcesByCountry(ctx context.Context, country *string) ([]map[string][]*str.WatchNowSource, *str.Response, error) {
	return w.fetchSources(ctx, fmt.Sprintf("watchnow/sources/%s", *country))
}

func (w *WatchNowService) fetchSources(ctx context.Context, url string) ([]map[string][]*str.WatchNowSource, *str.Response, error) {
	req, err := w.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []map[string][]*str.WatchNowSource{}
	resp, err := w.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
