// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
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

	printer.Println("fetch collection url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
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
	printer.Println("fetch history url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
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
	printer.Println("fetch watchlist url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetLastActivity Returns trakt user activity.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/last-activities/get-last-activity
func (s *SyncService) GetLastActivity(ctx context.Context) (*str.UserLastActivities, *str.Response, error) {
	var url string
	url = "sync/last_activities"

	printer.Println("fetch last activities url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.UserLastActivities)
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch activities err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPlaybackProgress Returns playback progress.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/playback/get-playback-progress
func (s *SyncService) GetPlaybackProgress(ctx context.Context, types *string, opts *uri.ListOptions) ([]*str.PlaybackProgress, *str.Response, error) {
	var url string
	if types != nil {
		url = fmt.Sprintf("sync/playback/%s", *types)
	} else {
		url = "sync/playback"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch playback url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.PlaybackProgress{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch playback err:" + err.Error())
		return nil, resp, err
	}
	return list, resp, nil
}

// RemovePlaybackItem removes playback item with selected id
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-playback/remove-a-playback-item
func (s *SyncService) RemovePlaybackItem(ctx context.Context, id *int) (*str.Response, error) {
	var url = fmt.Sprintf("sync/playback/%d", *id)
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf(consts.PlaybackNotFoundWithID, *id)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AddItemsToCollection add items to user's collection
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/add-to-collection/add-items-to-collection
func (s *SyncService) AddItemsToCollection(ctx context.Context, items *str.CollectionItems) (*str.CollectionAddResult, error) {
	var url = "sync/collection"
	printer.Println("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.CollectionAddResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetCollectedSeasons dedicated function do prepare collection: seasons format
func (s *SyncService) GetCollectedSeasons(ctx context.Context, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	// fetch collected shows
	strType := consts.Shows
	shows, resp, err := s.GetCollection(ctx, &strType, options)
	if err != nil {
		return nil, resp, err
	}
	collected := []str.Season{}
	for _, val := range shows {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		seasonsNumbers := []int{}
		for _, sitem := range *val.Seasons {
			seasonsNumbers = append(seasonsNumbers, *sitem.Number)
		}

		seasons, _, err := s.client.Shows.GetAllSeasonsForShow(ctx, val.Show.IDs.Slug, options)
		if err != nil {
			return nil, resp, err
		}
		for _, sitem := range seasons {
			if slices.Contains(seasonsNumbers, *sitem.Number) {
				s := str.Season{}
				s.IDs = sitem.IDs
				collected = append(collected, s)
			}
		}
	}

	strType = consts.Season
	list := []*str.ExportlistItem{}
	for _, citem := range collected {
		item := &str.ExportlistItem{}
		item.Type = &strType
		item.Season = &citem
		list = append(list, item)
	}

	return list, nil, nil
}
