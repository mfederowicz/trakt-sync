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

// GetSingleEpisodeForShow Returns a single episode's details.
//
// API docs: https://trakt.docs.apiary.io/#reference/episodes/summary/get-a-single-episode-for-a-show
func (s *EpisodesService) GetSingleEpisodeForShow(ctx context.Context, id *string, season *int, episode *int) (*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons/%d/episodes/%d", *id, *season, *episode)
	printer.Println("fetch single episode url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Episode)
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch episode err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetTrendingShows Returns the most watched shows over the last 24 hours.
// Shows with the most watchers are returned first.
// API docs: https://trakt.docs.apiary.io/#reference/shows/trending/get-trending-shows
func (s *ShowsService) GetTrendingShows(ctx context.Context, opts *uri.ListOptions) ([]*str.ShowsItem, *str.Response, error) {
	var url = "shows/trending"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPopularShows Returns the most popular shows.
// Popularity is calculated using the rating percentage and the number of ratings.
// API docs: https://trakt.docs.apiary.io/#reference/shows/popular/get-popular-shows
func (s *ShowsService) GetPopularShows(ctx context.Context, opts *uri.ListOptions) ([]*str.Show, *str.Response, error) {
	var url = "shows/popular"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Show{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetFavoritedShows Returns the most favorited shows in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/shows/favorited/get-favorited-shows
func (s *ShowsService) GetFavoritedShows(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.ShowsItem, *str.Response, error) {
	var url = fmt.Sprintf("shows/favorited/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPlayedShows Returns the most played (a single user can watch multiple episode multiple times) shows in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/shows/played/get-the-most-played-shows
func (s *ShowsService) GetPlayedShows(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.ShowsItem, *str.Response, error) {
	var url = fmt.Sprintf("shows/played/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}
	return list, resp, nil
}

// GetWatchedShows  Returns the most watched (unique users) shows in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/shows/watched/get-the-most-watched-shows
func (s *ShowsService) GetWatchedShows(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.ShowsItem, *str.Response, error) {
	var url = fmt.Sprintf("shows/watched/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetCollectedShows Returns the most collected (unique users) shows in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/shows/collected/get-the-most-collected-shows
func (s *ShowsService) GetCollectedShows(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.ShowsItem, *str.Response, error) {
	var url = fmt.Sprintf("shows/collected/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAnticipatedShows Returns the most anticipated shows based on the number of lists a show appears on.
// API docs: https://trakt.docs.apiary.io/#reference/shows/anticipated/get-the-most-anticipated-shows
func (s *ShowsService) GetAnticipatedShows(ctx context.Context, opts *uri.ListOptions) ([]*str.ShowsItem, *str.Response, error) {
	var url = "shows/anticipated"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch shows url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch shows err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
