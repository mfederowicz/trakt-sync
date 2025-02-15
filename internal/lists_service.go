// Package internal used for client and services
package internal

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// ListsService  handles communication with the lists related
// methods of the Trakt API.
type ListsService Service

// GetTrendingLists Returns all lists with the most likes and comments over the last 7 days.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/trending/get-trending-lists
func (l *ListsService) GetTrendingLists(ctx context.Context, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	var url string

	url = "lists/trending"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch trending url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.List{}
	resp, err := l.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPopularLists Returns the most popular lists. Popularity is calculated using total number of likes and comments..
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/popular/get-popular-lists 
func (l *ListsService) GetPopularLists(ctx context.Context, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	var url string

	url = "lists/popular"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch trending url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.List{}
	resp, err := l.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
