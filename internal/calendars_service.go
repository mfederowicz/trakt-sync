// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// CalendarsService  handles communication with the calendars related
// methods of the Trakt API.
type CalendarsService Service

// GetDVDReleases Returns all movies with a DVD release date during the time period specified.
//
// API docs: https://trakt.docs.apiary.io/#reference/calendars/my-dvd/get-dvd-releases
// API docs: https://trakt.docs.apiary.io/#reference/calendars/all-dvd/get-dvd-releases
func (c *CalendarsService) GetDVDReleases(ctx context.Context, actionType *string, startDate *string, days *int, opts *uri.ListOptions) ([]*str.CalendarList, *str.Response, error) {

	var urlStr = fmt.Sprintf("calendars/%s/dvd/%s/%d", *actionType, *startDate, *days)
	url, err := uri.AddQuery(urlStr, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch dvd calendars url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CalendarList{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch dvd calendars err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil

}

// GetMovies Returns all movies with a release date during the time period specified.
//
// API docs: https://trakt.docs.apiary.io/#reference/calendars/my-movies/get-movies
// API docs: https://trakt.docs.apiary.io/#reference/calendars/all-movies/get-movies
func (c *CalendarsService) GetMovies(ctx context.Context, actionType *string, startDate *string, days *int, opts *uri.ListOptions) ([]*str.CalendarList, *str.Response, error) {

	var url = fmt.Sprintf("calendars/%s/movies/%s/%d", *actionType, *startDate, *days)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch movies calendars url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CalendarList{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch movies calendars err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil

}

// GetSeasonPremieres Returns all show premieres (mid_season_premiere, season_premiere, series_premiere) airing during the time period specified.
//
// API docs: https://trakt.docs.apiary.io/#reference/calendars/my-season-premieres/get-season-premieres
// API docs: https://trakt.docs.apiary.io/#reference/calendars/all-season-premieres/get-season-premieres
func (c *CalendarsService) GetSeasonPremieres(ctx context.Context, actionType *string, startDate *string, days *int, opts *uri.ListOptions) ([]*str.CalendarList, *str.Response, error) {

	var url = fmt.Sprintf("calendars/%s/shows/premieres/%s/%d", *actionType, *startDate, *days)

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch season premieres url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CalendarList{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch season premieres err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil

}

// GetShows Returns all shows airing during the time period specified.
//
// API docs: https://trakt.docs.apiary.io/#reference/calendars/my-shows/get-shows
// API docs: https://trakt.docs.apiary.io/#reference/calendars/all-shows/get-shows
func (c *CalendarsService) GetShows(ctx context.Context, actionType *string, startDate *string, days *int, opts *uri.ListOptions) ([]*str.CalendarList, *str.Response, error) {

	var url = fmt.Sprintf("calendars/%s/shows/%s/%d", *actionType, *startDate, *days)

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch shows calendars url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CalendarList{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch shows calendars err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil

}

// GetNewShows Returns all new show premieres airing during the time period specified.
//
// API docs: https://trakt.docs.apiary.io/#reference/calendars/my-new-shows/get-new-shows
// API docs: https://trakt.docs.apiary.io/#reference/calendars/all-new-shows/get-new-shows
func (c *CalendarsService) GetNewShows(ctx context.Context, actionType *string, startDate *string, days *int, opts *uri.ListOptions) ([]*str.CalendarList, *str.Response, error) {

	var url = fmt.Sprintf("calendars/%s/shows/new/%s/%d", *actionType, *startDate, *days)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch new shows calendars url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CalendarList{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch new shows calendars err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil

}

// GetFinales Returns all show finales (mid_season_finale, season_finale, series_finale) airing during the time period specified.
//
// API docs: https://trakt.docs.apiary.io/#reference/calendars/my-finales/get-finales
// API docs: https://trakt.docs.apiary.io/#reference/calendars/all-finales/get-finales
func (c *CalendarsService) GetFinales(ctx context.Context, actionType *string, startDate *string, days *int, opts *uri.ListOptions) ([]*str.CalendarList, *str.Response, error) {

	var url = fmt.Sprintf("calendars/%s/shows/finales/%s/%d", *actionType, *startDate, *days)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("fetch finales calendars url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CalendarList{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		fmt.Println("fetch finales calendars err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil

}
