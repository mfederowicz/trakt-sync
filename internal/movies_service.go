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
func (m *MoviesService) GetMovie(ctx context.Context, id *string, opts *uri.ListOptions) (*str.Movie, *str.Response, error) {
	var url = fmt.Sprintf("movies/%s", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movie url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	movie := new(str.Movie)
	resp, err := m.client.Do(ctx, req, &movie)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found movie for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch movie err:" + err.Error())
		return nil, resp, err
	}

	return movie, resp, nil
}

// GetTrendingMovies Returns the most watched movies over the last 24 hours.
// Movies with the most watchers are returned first.
// API docs: https://trakt.docs.apiary.io/#reference/movies/trending/get-trending-movies
func (m *MoviesService) GetTrendingMovies(ctx context.Context, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
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

	list := []*str.MoviesItem{}
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
func (m *MoviesService) GetFavoritedMovies(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.MoviesItem, *str.Response, error) {
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

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPlayedMovies Returns the most played (a single user can watch multiple times) movies in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/movies/played/get-the-most-played-movies
func (m *MoviesService) GetPlayedMovies(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.MoviesItem, *str.Response, error) {
	var url = fmt.Sprintf("movies/played/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchedMovies  Returns the most watched (unique users) movies in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/movies/watched/get-the-most-watched-movies
func (m *MoviesService) GetWatchedMovies(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.MoviesItem, *str.Response, error) {
	var url = fmt.Sprintf("movies/watched/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetCollectedMovies Returns the most collected (unique users) movies in the specified time period, defaulting to weekly.
// All stats are relative to the specific time period.
// API docs: https://trakt.docs.apiary.io/#reference/movies/collected/get-the-most-collected-movies
func (m *MoviesService) GetCollectedMovies(ctx context.Context, opts *uri.ListOptions, period *string) ([]*str.MoviesItem, *str.Response, error) {
	var url = fmt.Sprintf("movies/collected/%s", *period)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAnticipatedMovies Returns the most anticipated movies based on the number of lists a movie appears on.
// API docs: https://trakt.docs.apiary.io/#reference/movies/anticipated/get-the-most-anticipated-movies
func (m *MoviesService) GetAnticipatedMovies(ctx context.Context, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
	var url = "movies/anticipated"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetBoxoffice Returns the top 10 grossing movies in the U.S. box office last weekend. Updated every Monday morning.
// API docs: https://trakt.docs.apiary.io/#reference/movies/box-office/get-the-weekend-box-office
func (m *MoviesService) GetBoxoffice(ctx context.Context, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
	var url = "movies/boxoffice"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch movies url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch movies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetRecentlyUpdatedMovies Returns all movies updated since the specified UTC date and time.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/updates/get-recently-updated-movies
func (m *MoviesService) GetRecentlyUpdatedMovies(ctx context.Context, startDate *string, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("movies/updates/%s", *startDate)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updates url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.MoviesItem{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch updates err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetRecentlyUpdatedMoviesTraktIDs Returns all movie Trakt IDs updated since the specified UTC date and time.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/updated-ids/get-recently-updated-movie-trakt-ids
func (m *MoviesService) GetRecentlyUpdatedMoviesTraktIDs(ctx context.Context, startDate *string, opts *uri.ListOptions) ([]*int, *str.Response, error) {
	var url string

	url = fmt.Sprintf("movies/updates/id/%s", *startDate)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updates url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*int{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch updates err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllMovieAliases Returns all title aliases for a movie. Includes country where name is different.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/aliases/get-all-movie-aliases
func (m *MoviesService) GetAllMovieAliases(ctx context.Context, id *string) ([]*str.Alias, *str.Response, error) {
	url := fmt.Sprintf("movies/%s/aliases", *id)
	printer.Println("fetch aliases url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Alias{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch aliases err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllMovieReleases Returns all releases for a movie including country, certification, release date, release type, and note.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/releases/get-all-movie-releases
func (m *MoviesService) GetAllMovieReleases(ctx context.Context, id *string, country *string) ([]*str.Release, *str.Response, error) {
	var url string
	if country != nil {
		url = fmt.Sprintf("movies/%s/releases/%s", *id, *country)
	} else {
		url = fmt.Sprintf("movies/%s/releases", *id)
	}

	printer.Println("fetch releases url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Release{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch releases err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
