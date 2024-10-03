// Package internal used for client and services
package internal

import (
	"context"
	"fmt"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SearchService  handles communication with the search related
// methods of the Trakt API.
type SearchService Service


// GetTextQueryResults Search all text fields that a media object contains
// (i.e. title, overview, etc). Results are ordered by the most relevant score.
// Specify the type of results by sending a single value or a comma delimited string for multiple types.
//
// API docs: https://trakt.docs.apiary.io/#reference/search/text-query/get-text-query-results
func (s *SearchService) GetTextQueryResults(ctx context.Context, searchType *string, opts *uri.ListOptions) ([]*str.SearchListItem, *str.Response, error) {
	var url = fmt.Sprintf("search/%s", *searchType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch text search url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.SearchListItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch text search err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetIDLookupResults Lookup items by their Trakt, IMDB, TMDB, or TVDB ID.
// If you use the search url without a type it might return multiple items
// if the id_type is not globally unique. Specify the type of results by
// sending a single value or a comma delimited string for multiple types.
//
// API docs: https://trakt.docs.apiary.io/#reference/search/id-lookup/get-id-lookup-results
func (s *SearchService) GetIDLookupResults(ctx context.Context, formatType *string, id *string, opts *uri.ListOptions) ([]*str.SearchListItem, *str.Response, error) {
	var url = fmt.Sprintf("search/%s/%s", *formatType, *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch id lookup search url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.SearchListItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lookup search err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
