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
func (s *SyncService) GetWatchedHistory(ctx context.Context, id *int, types *string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if types != nil {
		url = fmt.Sprintf("sync/history/%s", *types)
	} else {
		url = "sync/history"
	}

	if *id > consts.ZeroValue {
		url = fmt.Sprintf(url+"/%d", *id)
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
func (s *SyncService) GetWatchlist(ctx context.Context, types *string, sortBy *string, sortHow *string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if types != nil && sortBy != nil && sortHow != nil {
		url = fmt.Sprintf("sync/watchlist/%s/%s/%s", *types, *sortBy, *sortHow)
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
func (s *SyncService) AddItemsToCollection(ctx context.Context, items *str.ItemsList) (*str.CollectionAddResult, error) {
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

// RemoveItemsFromCollection remove items from user's collection
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-from-collection/remove-items-from-collection
func (s *SyncService) RemoveItemsFromCollection(ctx context.Context, items *str.ItemsList) (*str.CollectionRemoveResult, error) {
	var url = "sync/collection/remove"
	printer.Println("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.CollectionRemoveResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetWatched Returns all movies or shows a user has watched sorted by most plays.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched
func (s *SyncService) GetWatched(ctx context.Context, watchType *string, opts *uri.ListOptions) ([]*str.UserWatched, *str.Response, error) {
	var url string
	url = fmt.Sprintf("sync/watched/%s", *watchType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("get watched url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	watched := []*str.UserWatched{}
	resp, err := s.client.Do(ctx, req, &watched)

	if err != nil {
		return nil, resp, err
	}

	return watched, resp, nil
}

// AddItemsToHistory add items to user's history
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/add-to-history/add-items-to-watched-history
func (s *SyncService) AddItemsToHistory(ctx context.Context, items *str.HistoryItems) (*str.AddResult, error) {
	var url = "sync/history"
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.AddResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// RemoveItemsFromHistory remove items from user's history
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-from-history/remove-items-from-history
func (s *SyncService) RemoveItemsFromHistory(ctx context.Context, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	var url = "sync/history/remove"
	printer.Println("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.RemoveResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetRatings Returns users ratings.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-ratings/get-ratings
func (s *SyncService) GetRatings(ctx context.Context, types *string, rating *string, opts *uri.ListOptions) ([]*str.RatingListItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("sync/ratings/%s", *types)
	if len(*rating) > consts.ZeroValue {
		url = fmt.Sprintf("sync/ratings/%s/%s", *types, *rating)
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch ratings url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.RatingListItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// RemoveItemsFromRatings Remove ratings for one or more items.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-ratings/remove-ratings
func (s *SyncService) RemoveItemsFromRatings(ctx context.Context, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	var url = "sync/ratings/remove"
	printer.Println("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.RemoveResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// AddItemsToRatings Rate one or more items. Accepts shows, seasons, episodes and movies.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/add-ratings/add-new-ratings
func (s *SyncService) AddItemsToRatings(ctx context.Context, items *str.RatingItems) (*str.AddResult, error) {
	var url = "sync/ratings"
	printer.Println("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.AddResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// UpdateWatchlist Update the watchlist by sending 1 or more parameters.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-watchlist/update-watchlist
func (s *SyncService) UpdateWatchlist(ctx context.Context, update *str.PersonalList) (*str.PersonalList, error) {
	var url = "sync/watchlist"
	printer.Println("update watchlist")
	req, err := s.client.NewRequest(http.MethodPut, url, update)
	if err != nil {
		return nil, err
	}

	result := new(str.PersonalList)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// RemoveItemsFromWatchlist Remove one or more items from a user's watchlist.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-from-watchlist/remove-items-from-watchlist
func (s *SyncService) RemoveItemsFromWatchlist(context context.Context, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	var url = "sync/watchlist/remove"
	printer.Println("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.RemoveResult)
	_, err = s.client.Do(context, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// AddItemsToWatchlist Add one of more items to a user's watchlist.
// Accepts shows, seasons, episodes and movies. If only a show is passed,
// only the show itself will be added. If seasons are specified, all of
// those seasons will be added.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-watchlist/add-items-to-watchlist
func (s *SyncService) AddItemsToWatchlist(ctx context.Context, items *str.HistoryItems) (*str.AddResult, error) {
	var url = "sync/watchlist"
	printer.Println("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.AddResult)
	_, err = s.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
