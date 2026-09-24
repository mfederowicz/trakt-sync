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

// SearchService  handles communication with the search related
// methods of the Trakt API.
type SearchService Service

// GetTextQueryResults Search all text fields that a media object contains
// (i.e. title, overview, etc). Results are ordered by the most relevant score.
// Specify the type of results by sending a single value or a comma delimited string for multiple types.
//
// API docs: https://docs.trakt.tv/reference/getsearchquery
func (s *SearchService) GetTextQueryResults(ctx context.Context, searchType *string, opts *uri.ListOptions) ([]*str.SearchListItem, *str.Response, error) {
	var url = fmt.Sprintf("search/%s", *searchType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch text search url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
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
// API docs: https://docs.trakt.tv/reference/getsearchlookup
func (s *SearchService) GetIDLookupResults(ctx context.Context, formatType *string, id *string, opts *uri.ListOptions) ([]*str.SearchListItem, *str.Response, error) {
	var url = fmt.Sprintf("search/%s/%s", *formatType, *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch id lookup search url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
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

// GetExactTextQueryResults Search for exact movie or show matches for the query.
//
// API docs: https://docs.trakt.tv/reference/getsearchexact
func (s *SearchService) GetExactTextQueryResults(ctx context.Context, searchType *string, opts *uri.ListOptions) ([]*str.SearchListItem, *str.Response, error) {
	var url = fmt.Sprintf("search/%s/exact", *searchType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.SearchListItem{}
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetTrendingSearches Get globally trending recent searches by type.
//
// API docs: https://docs.trakt.tv/reference/getsearchtrending
func (s *SearchService) GetTrendingSearches(ctx context.Context, searchType *string, opts *uri.ListOptions) ([]*str.SearchTrendingItem, *str.Response, error) {
	var url = fmt.Sprintf("search/recent_by_id/global/%s", *searchType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.SearchTrendingItem{}
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// AddRecentSearch Add a recent search to the global search trends.
//
// API docs: https://docs.trakt.tv/reference/postsearchrecentadd
func (s *SearchService) AddRecentSearch(ctx context.Context, search *str.RecentSearch) (*str.Response, error) {
	var url = "search/recent"
	req, err := s.client.NewRequest(http.MethodPost, url, search)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveRecentSearch Remove a recent search from the global search trends.
//
// API docs: https://docs.trakt.tv/reference/postsearchrecentremove
func (s *SearchService) RemoveRecentSearch(ctx context.Context, search *str.RecentSearch) (*str.Response, error) {
	var url = "search/recent/remove"
	req, err := s.client.NewRequest(http.MethodPost, url, search)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
