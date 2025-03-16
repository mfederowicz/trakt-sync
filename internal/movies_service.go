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

// MoviesService  handles communication with the movies related
// methods of the Trakt API.
type MoviesService Service

// GetMovie Returns movie object.
func (m *MoviesService) GetMovie(ctx context.Context, id *int) (*str.Movie, *str.Response, error) {
	var url = fmt.Sprintf("movies/%d", *id)
	printer.Println("fetch movie url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	movie := new(str.Movie)
	resp, err := m.client.Do(ctx, req, &movie)

	if err != nil {
		printer.Println("fetch movie err:" + err.Error())
		return nil, resp, err
	}

	return movie, resp, nil
}

// GetTrendingMovies Returns the most watched movies over the last 24 hours.
// Movies with the most watchers are returned first.
// API docs: https://trakt.docs.apiary.io/#reference/movies/trending/get-trending-movies
func (m *MoviesService) GetTrendingMovies(ctx context.Context, opts *uri.ListOptions) ([]*str.TrendingMovie, *str.Response, error) {
	var url = "movies/trending"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.TrendingMovie{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPopularMovies Returns the most popular movies. 
// Popularity is calculated using the rating percentage and the number of ratings.
// API docs: https://trakt.docs.apiary.io/#reference/movies/popular/get-popular-movies
func (m *MoviesService) GetPopularMovies(ctx context.Context, opts *uri.ListOptions) ([]*str.Movie, *str.Response, error) {
	var url = "movies/popular"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Movie{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetFavoritedMovies Returns the most favorited movies in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/movies/favorited/get-favorited-movies
func (m *MoviesService) GetFavoritedMovies(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.FavoritedMovie, *str.Response, error) {
	var url = fmt.Sprintf("movies/favorited/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.FavoritedMovie{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
