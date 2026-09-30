package trakt

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SyncService  handles communication with the sync related
// methods of the Trakt API.
type SyncService Service

// GetCollection Get all collected items in a user's collection. An empty types returns all types.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-collection/get-collection
func (s *SyncService) GetCollection(ctx context.Context, types string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if types != consts.EmptyString {
		url = fmt.Sprintf("sync/collection/%s", types)
	} else {
		url = "sync/collection"
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	s.client.debug("fetch collection url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		s.client.debug("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchedHistory Returns movies and episodes that a user has watched, sorted by most recent.
// An empty types returns all types; an id of 0 returns all entries instead of one history item.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched-history
func (s *SyncService) GetWatchedHistory(ctx context.Context, id int, types string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if types != consts.EmptyString {
		url = fmt.Sprintf("sync/history/%s", types)
	} else {
		url = "sync/history"
	}

	if id > consts.ZeroValue {
		url = fmt.Sprintf(url+"/%d", id)
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	s.client.debug("fetch history url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		s.client.debug("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchlist Returns all items in a user's watchlist filtered by type.
// The type and sort segments are sent only when types, sortBy and sortHow are all set.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-watchlist/get-watchlist
func (s *SyncService) GetWatchlist(ctx context.Context, types string, sortBy string, sortHow string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if types != consts.EmptyString && sortBy != consts.EmptyString && sortHow != consts.EmptyString {
		url = fmt.Sprintf("sync/watchlist/%s/%s/%s", types, sortBy, sortHow)
	} else {
		url = "sync/watchlist"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	s.client.debug("fetch watchlist url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		s.client.debug("fetch lists err:" + err.Error())
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

	s.client.debug("fetch last activities url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.UserLastActivities)
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		s.client.debug("fetch activities err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPlaybackProgress Returns playback progress; types movies or episodes narrows it, an empty types returns both.
//
// API docs: https://docs.trakt.tv/reference/getsyncprogressplayback
func (s *SyncService) GetPlaybackProgress(ctx context.Context, types string, opts *uri.ListOptions) ([]*str.PlaybackProgress, *str.Response, error) {
	var url string
	if types != consts.EmptyString {
		url = fmt.Sprintf("sync/playback/%s", types)
	} else {
		url = "sync/playback"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	s.client.debug("fetch playback url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.PlaybackProgress{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		s.client.debug("fetch playback err:" + err.Error())
		return nil, resp, err
	}
	return list, resp, nil
}

// RemovePlaybackItem removes playback item with selected id
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-playback/remove-a-playback-item
func (s *SyncService) RemovePlaybackItem(ctx context.Context, id int) (*str.Response, error) {
	var url = fmt.Sprintf("sync/playback/%d", id)
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf(consts.PlaybackNotFoundWithID, id)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AddItemsToCollection add items to user's collection
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/add-to-collection/add-items-to-collection
func (s *SyncService) AddItemsToCollection(ctx context.Context, items *str.ItemsList) (*str.CollectionAddResult, *str.Response, error) {
	var url = "sync/collection"
	s.client.debug("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.CollectionAddResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// GetCollectedSeasons dedicated function do prepare collection: seasons format
func (s *SyncService) GetCollectedSeasons(ctx context.Context, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	// fetch collected shows
	strType := consts.Shows
	shows, resp, err := s.GetCollection(ctx, strType, options)
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
func (s *SyncService) RemoveItemsFromCollection(ctx context.Context, items *str.ItemsList) (*str.CollectionRemoveResult, *str.Response, error) {
	var url = "sync/collection/remove"
	s.client.debug("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.CollectionRemoveResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// GetWatched Returns all movies or shows a user has watched sorted by most plays.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched
func (s *SyncService) GetWatched(ctx context.Context, watchType string, opts *uri.ListOptions) ([]*str.UserWatched, *str.Response, error) {
	var url string
	url = fmt.Sprintf("sync/watched/%s", watchType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	s.client.debug("get watched url:" + url)
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
func (s *SyncService) AddItemsToHistory(ctx context.Context, items *str.HistoryItems) (*str.AddResult, *str.Response, error) {
	var url = "sync/history"
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.AddResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// RemoveItemsFromHistory remove items from user's history
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-from-history/remove-items-from-history
func (s *SyncService) RemoveItemsFromHistory(ctx context.Context, items *str.ItemsToRemove) (*str.RemoveResult, *str.Response, error) {
	var url = "sync/history/remove"
	s.client.debug("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.RemoveResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// GetRatings Returns users ratings. An empty rating returns all ratings.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-ratings/get-ratings
func (s *SyncService) GetRatings(ctx context.Context, types string, rating string, opts *uri.ListOptions) ([]*str.RatingListItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("sync/ratings/%s", types)
	if len(rating) > consts.ZeroValue {
		url = fmt.Sprintf("sync/ratings/%s/%s", types, rating)
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	s.client.debug("fetch ratings url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.RatingListItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		s.client.debug("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// RemoveItemsFromRatings Remove ratings for one or more items.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-ratings/remove-ratings
func (s *SyncService) RemoveItemsFromRatings(ctx context.Context, items *str.ItemsToRemove) (*str.RemoveResult, *str.Response, error) {
	var url = "sync/ratings/remove"
	s.client.debug("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.RemoveResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// AddItemsToRatings Rate one or more items. Accepts shows, seasons, episodes and movies.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/add-ratings/add-new-ratings
func (s *SyncService) AddItemsToRatings(ctx context.Context, items *str.RatingItems) (*str.AddResult, *str.Response, error) {
	var url = "sync/ratings"
	s.client.debug("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.AddResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// UpdateWatchlist Update the watchlist by sending 1 or more parameters.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-watchlist/update-watchlist
func (s *SyncService) UpdateWatchlist(ctx context.Context, update *str.PersonalList) (*str.PersonalList, *str.Response, error) {
	var url = "sync/watchlist"
	s.client.debug("update watchlist")
	req, err := s.client.NewRequest(http.MethodPut, url, update)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.PersonalList)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// UpdateFavorites Update the favorites list by sending 1 or more parameters.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-favorites/update-favorites
func (s *SyncService) UpdateFavorites(ctx context.Context, update *str.PersonalList) (*str.PersonalList, *str.Response, error) {
	var url = "sync/favorites"
	s.client.debug("update favorites")
	req, err := s.client.NewRequest(http.MethodPut, url, update)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.PersonalList)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// UpdateWatchlistItem Update the notes on a single watchlist item.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-watchlist-item/update-a-watchlist-item
func (s *SyncService) UpdateWatchlistItem(ctx context.Context, itemID int, update *str.WatchlistItem) (*str.Response, error) {
	var url string

	url = fmt.Sprintf("sync/watchlist/%d", itemID)
	s.client.debug("update notes")
	req, err := s.client.NewRequest(http.MethodPut, url, update)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveItemsFromWatchlist Remove one or more items from a user's watchlist.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/remove-from-watchlist/remove-items-from-watchlist
func (s *SyncService) RemoveItemsFromWatchlist(ctx context.Context, items *str.ItemsToRemove) (*str.RemoveResult, *str.Response, error) {
	var url = "sync/watchlist/remove"
	s.client.debug("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.RemoveResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// AddItemsToWatchlist Add one of more items to a user's watchlist.
// Accepts shows, seasons, episodes and movies. If only a show is passed,
// only the show itself will be added. If seasons are specified, all of
// those seasons will be added.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-watchlist/add-items-to-watchlist
func (s *SyncService) AddItemsToWatchlist(ctx context.Context, items *str.HistoryItems) (*str.AddResult, *str.Response, error) {
	var url = "sync/watchlist"
	s.client.debug("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.AddResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// ReorderWatchlistItems Reorder all items on a user's watchlist by sending the updated rank of list item ids.
// Use the /sync/watchlist method to get all list item ids.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/reorder-watchlist/reorder-watchlist-items
func (s *SyncService) ReorderWatchlistItems(ctx context.Context, reorder *str.ItemsToReorder) (*str.ReorderResults, *str.Response, error) {
	var url = "sync/watchlist/reorder"
	s.client.debug("reorder watchlist")
	req, err := s.client.NewRequest(http.MethodPost, url, reorder)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.ReorderResults)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// GetFavorites Returns all items in a user's favorites filtered by type.
// The type and sort segments are sent only when types, sortBy and sortHow are all set.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/get-favorites/get-favorites
func (s *SyncService) GetFavorites(ctx context.Context, types string, sortBy string, sortHow string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if types != consts.EmptyString && sortBy != consts.EmptyString && sortHow != consts.EmptyString {
		url = fmt.Sprintf("sync/favorites/%s/%s/%s", types, sortBy, sortHow)
	} else {
		url = "sync/favorites"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	s.client.debug("fetch favorites url:" + url)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		s.client.debug("fetch favorites err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// AddItemsToFavorites add items to favorites.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/update-favorites/add-items-to-favorites
func (s *SyncService) AddItemsToFavorites(ctx context.Context, items *str.HistoryItems) (*str.AddResult, *str.Response, error) {
	var url = "sync/favorites"
	s.client.debug("add items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.AddResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// RemoveItemsFromFavorites remove items from favorites.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/remove-from-favorites/remove-items-from-favorites
func (s *SyncService) RemoveItemsFromFavorites(ctx context.Context, items *str.ItemsToRemove) (*str.RemoveResult, *str.Response, error) {
	var url = "sync/favorites/remove"
	s.client.debug("remove items")
	req, err := s.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.RemoveResult)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// ReorderFavoritesItems Reorder all items on a user's favorites by sending the updated rank of list item ids.
// Use the /sync/favorites method to get all list item ids.
//
// API docs:https://trakt.docs.apiary.io/#reference/sync/reorder-favorites/reorder-favorited-items
func (s *SyncService) ReorderFavoritesItems(ctx context.Context, reorder *str.ItemsToReorder) (*str.ReorderResults, *str.Response, error) {
	var url = "sync/favorites/reorder"
	s.client.debug("reorder favorites")
	req, err := s.client.NewRequest(http.MethodPost, url, reorder)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.ReorderResults)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// UpdateFavoriteItem Update the notes on a single favorite item.
//
// API docs: https://trakt.docs.apiary.io/#reference/sync/update-favorite-item/update-a-favorite-item
func (s *SyncService) UpdateFavoriteItem(ctx context.Context, itemID int, update *str.FavoriteItem) (*str.Response, error) {
	var url string

	url = fmt.Sprintf("sync/favorites/%d", itemID)
	s.client.debug("update notes")
	req, err := s.client.NewRequest(http.MethodPut, url, update)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetMinimalCollection Returns the movie or episode collection in a minimal format: Trakt ID -> collected_at.
//
// API docs: https://docs.trakt.tv/reference/getsynccollectionminimalmovies
// API docs: https://docs.trakt.tv/reference/getsynccollectionminimalepisodes
func (s *SyncService) GetMinimalCollection(ctx context.Context, strType string, opts *uri.ListOptions) (str.MinimalCollection, *str.Response, error) {
	var url = fmt.Sprintf("sync/collection/minimal/%s", strType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := str.MinimalCollection{}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetMinimalShowCollection Returns the show collection in a minimal format: show Trakt ID -> season -> episode -> collected_at.
//
// API docs: https://docs.trakt.tv/reference/getsynccollectionminimalshows
func (s *SyncService) GetMinimalShowCollection(ctx context.Context, opts *uri.ListOptions) (str.MinimalShowCollection, *str.Response, error) {
	var url = "sync/collection/minimal/shows"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := str.MinimalShowCollection{}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetUpNext Returns the up next progress: shows with their next episode to watch.
//
// API docs: https://docs.trakt.tv/reference/getsyncprogressupnextstandard
func (s *SyncService) GetUpNext(ctx context.Context, opts *uri.SyncProgressOptions) ([]*str.ShowProgress, *str.Response, error) {
	return s.getShowProgress(ctx, "sync/progress/up_next", opts)
}

// GetWatchedProgress Returns the watched progress of the user's shows.
//
// API docs: https://docs.trakt.tv/reference/getsyncprogresswatched
func (s *SyncService) GetWatchedProgress(ctx context.Context, opts *uri.SyncProgressOptions) ([]*str.ShowProgress, *str.Response, error) {
	return s.getShowProgress(ctx, "sync/progress/watched", opts)
}

// GetUpNextNitro Returns the up next progress for intent-based clients, with media filters.
//
// API docs: https://docs.trakt.tv/reference/getsyncprogressupnextnitro
func (s *SyncService) GetUpNextNitro(ctx context.Context, opts *uri.UpNextNitroOptions) ([]*str.ShowProgress, *str.Response, error) {
	return s.getShowProgress(ctx, "sync/progress/up_next_nitro", opts)
}

func (s *SyncService) getShowProgress(ctx context.Context, url string, opts any) ([]*str.ShowProgress, *str.Response, error) {
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowProgress{}
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
