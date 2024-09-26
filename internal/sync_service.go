// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SyncService  handles communication with the sync related
// methods of the Trakt API.
type SyncService Service

// GetCollection Get all collected items in a user's collection.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-collection/get-collection
func (s *SyncService) GetCollection(ctx context.Context, types *string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {

	var url string

	if types != nil {
		url = fmt.Sprintf("sync/collection/%s", *types)
	} else {
		url = "sync/collection"
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch collection url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchedHistory Returns movies and episodes that a user has watched, sorted by most recent.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched-history
func (s *SyncService) GetWatchedHistory(ctx context.Context, types *string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {

	var url string

	if types != nil {
		url = fmt.Sprintf("sync/history/%s", *types)
	} else {
		url = "sync/history"
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("fetch history url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchlist Returns all items in a user's watchlist filtered by type.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-watchlist/get-watchlist
func (s *SyncService) GetWatchlist(ctx context.Context, types *string, sort *string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {

	var url string

	if types != nil && sort != nil {
		url = fmt.Sprintf("sync/watchlist/%s/%s", *types, *sort)
	} else {
		url = "sync/watchlist"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("fetch watchlist url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
