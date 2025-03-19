// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// ShowsService  handles communication with the shows related
// methods of the Trakt API.
type ShowsService Service

// GetShow Returns episode object.
func (s *ShowsService) GetShow(ctx context.Context, id *string, opts *uri.ListOptions) (*str.Show, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch show url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Show)
	resp, err := s.client.Do(ctx, req, &result)
	
	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("show not found with traktId:%s", *id)
	}

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
