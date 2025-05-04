// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found aliases for id/slug:%s", *id)
	}

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

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found releases for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch releases err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllMovieTranslations Returns all translations for a movie, including language and translated values for title, tagline and overview.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/translations/get-all-movie-translations
func (m *MoviesService) GetAllMovieTranslations(ctx context.Context, id *string, language *string) ([]*str.Translation, *str.Response, error) {
	var url string
	if *language != consts.EmptyString {
		url = fmt.Sprintf("movies/%s/translations/%s", *id, *language)
	} else {
		url = fmt.Sprintf("movies/%s/translations", *id)
	}

	printer.Println("fetch translations url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Translation{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found translations for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch translations err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllMovieComments Returns all top level comments for a movie.
// By default, the newest comments are returned first.
// Other sorting options include oldest, most likes, most replies, highest rated, lowest rated, and most plays..
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/comments/get-all-movie-comments
func (m *MoviesService) GetAllMovieComments(ctx context.Context, id *string, sort *string, opts *uri.ListOptions) ([]*str.Comment, *str.Response, error) {
	var url string
	if *sort != consts.EmptyString {
		url = fmt.Sprintf("movies/%s/comments/%s", *id, *sort)
	} else {
		url = fmt.Sprintf("movies/%s/comments", *id)
	}

	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch comments url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Comment{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found comments for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch comments err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetListsContainingMovie Returns all lists that contain this movie.
// By default, personal lists are returned sorted by the most popular.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/lists/get-lists-containing-this-movie
func (m *MoviesService) GetListsContainingMovie(ctx context.Context, id *string, t *string, sort *string, opts *uri.ListOptions) ([]*str.PersonalList, *str.Response, error) {
	var url string
	if *t != consts.EmptyString && *sort != consts.EmptyString {
		url = fmt.Sprintf("movies/%s/lists/%s/%s", *id, *t, *sort)
	} else {
		url = fmt.Sprintf("movies/%s/lists", *id)
	}

	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch lists url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.PersonalList{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found lists for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllPeopleForMovie Returns all cast and crew for a movie.
// Each cast member will have a characters array and a standard person object.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/people/get-all-people-for-a-movie
func (m *MoviesService) GetAllPeopleForMovie(ctx context.Context, id *string, opts *uri.ListOptions) (*str.MoviePeople, *str.Response, error) {
	var url string

	url = fmt.Sprintf("movies/%s/people", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch people url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.MoviePeople)
	resp, err := m.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found people for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch people err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetMovieRatings Returns rating (between 0 and 10) and distribution for a movie.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/ratings/get-movie-ratings
func (m *MoviesService) GetMovieRatings(ctx context.Context, id *string) (*str.MovieRatings, *str.Response, error) {
	url := fmt.Sprintf("movies/%s/ratings", *id)
	printer.Println("fetch ratings url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.MovieRatings)
	resp, err := m.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found ratings for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch ratings err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetRelatedMovies Returns related and similar movies.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/related/get-related-movies
func (m *MoviesService) GetRelatedMovies(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.Movie, *str.Response, error) {
	var url string
	url = fmt.Sprintf("movies/%s/related", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch related url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Movie{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found related for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch related err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetMovieStats Returns lots of movie stats.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/stats/get-movie-stats
func (m *MoviesService) GetMovieStats(ctx context.Context, id *string) (*str.MovieStats, *str.Response, error) {
	url := fmt.Sprintf("movies/%s/stats", *id)
	printer.Println("fetch stats url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.MovieStats)
	resp, err := m.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found stats for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch stats err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetMovieStudios Returns all studios for movie.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/studios/get-movie-studios
func (m *MoviesService) GetMovieStudios(ctx context.Context, id *string) ([]*str.Studio, *str.Response, error) {
	var url = fmt.Sprintf("movies/%s/studios", *id)
	printer.Println("fetch studios url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Studio{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found studios for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch studios err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetMovieWatching Returns all users watching this movie right now.
//
// API docs:  https://trakt.docs.apiary.io/#reference/movies/studios/get-users-watching-right-now
func (m *MoviesService) GetMovieWatching(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.UserProfile, *str.Response, error) {
	var url = fmt.Sprintf("movies/%s/watching", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch watching url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.UserProfile{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found watching for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch watching err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetMovieVideos Returns all videos including trailers, teasers, clips, and featurettes.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/videos/get-all-videos
func (m *MoviesService) GetMovieVideos(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.Video, *str.Response, error) {
	var url = fmt.Sprintf("movies/%s/videos", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch video url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Video{}
	resp, err := m.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found video for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch video err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// RefreshMovieMetadata Queue this movie for a full metadata and image refresh.
// It might take up to 8 hours for the updated metadata to be availabe through the API.
//
// API docs: https://trakt.docs.apiary.io/#reference/movies/refresh/refresh-movie-metadata
func (m *MoviesService) RefreshMovieMetadata(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("movies/%s/refresh", *id)
	printer.Println("refresh movie:" + url)
	req, err := m.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
