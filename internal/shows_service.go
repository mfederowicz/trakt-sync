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

// GetRecentlyUpdatedShows Returns all shows updated since the specified UTC date and time.
// API docs: https://trakt.docs.apiary.io/#reference/shows/updates/get-recently-updated-shows
func (s *ShowsService) GetRecentlyUpdatedShows(ctx context.Context, startDate *string, opts *uri.ListOptions) ([]*str.ShowsItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("shows/updates/%s", *startDate)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updates url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ShowsItem{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch updates err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetRecentlyUpdatedShowsTraktIDs Returns all show Trakt IDs updated since the specified UTC date and time.
// API docs: https://trakt.docs.apiary.io/#reference/shows/updated-ids/get-recently-updated-show-trakt-ids
func (s *ShowsService) GetRecentlyUpdatedShowsTraktIDs(ctx context.Context, startDate *string, opts *uri.ListOptions) ([]*int, *str.Response, error) {
	var url string

	url = fmt.Sprintf("shows/updates/id/%s", *startDate)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updates url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*int{}
	resp, err := s.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch updates err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllShowAliases Returns all title aliases for a show. Includes country where name is different.
// API docs: https://trakt.docs.apiary.io/#reference/shows/aliases/get-all-show-aliases
func (s *ShowsService) GetAllShowAliases(ctx context.Context, id *string) ([]*str.Alias, *str.Response, error) {
	url := fmt.Sprintf("shows/%s/aliases", *id)
	printer.Println("fetch aliases url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Alias{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found aliases for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch aliases err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllShowCertifications Returns all content certifications for a show, including the country.
// API docs: https://trakt.docs.apiary.io/#reference/shows/certifications/get-all-show-certifications
func (s *ShowsService) GetAllShowCertifications(ctx context.Context, id *string) ([]*str.Certification, *str.Response, error) {
	url := fmt.Sprintf("shows/%s/certifications", *id)
	printer.Println("fetch certifications url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Certification{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found certifications for id/slug:%s", *id)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, nil, fmt.Errorf("fetch certifications: internal server error for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch certifications err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllShowTranslations Returns all translations for a show, including language and translated values for title, tagline and overview.
// API docs: https://trakt.docs.apiary.io/#reference/shows/translations/get-all-show-translations
func (s *ShowsService) GetAllShowTranslations(ctx context.Context, id *string, language *string) ([]*str.Translation, *str.Response, error) {
	var url string
	if *language != consts.EmptyString {
		url = fmt.Sprintf("shows/%s/translations/%s", *id, *language)
	} else {
		url = fmt.Sprintf("shows/%s/translations", *id)
	}

	printer.Println("fetch translations url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Translation{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found translations for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch translations err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllShowComments Returns all top level comments for a show.
// By default, the newest comments are returned first.
// Other sorting options include oldest, most likes, most replies, highest rated, lowest rated, and most plays..
// API docs: https://trakt.docs.apiary.io/#reference/shows/comments/get-all-show-comments
func (s *ShowsService) GetAllShowComments(ctx context.Context, id *string, sort *string, opts *uri.ListOptions) ([]*str.Comment, *str.Response, error) {
	var url string
	if *sort != consts.EmptyString {
		url = fmt.Sprintf("shows/%s/comments/%s", *id, *sort)
	} else {
		url = fmt.Sprintf("shows/%s/comments", *id)
	}

	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch comments url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Comment{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found comments for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch comments err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetListsContainingShow Returns all lists that contain this show.
// By default, personal lists are returned sorted by the most popular.
// API docs: https://trakt.docs.apiary.io/#reference/shows/lists/get-lists-containing-this-show
func (s *ShowsService) GetListsContainingShow(ctx context.Context, id *string, t *string, sort *string, opts *uri.ListOptions) ([]*str.PersonalList, *str.Response, error) {
	var url string
	if *t != consts.EmptyString && *sort != consts.EmptyString {
		url = fmt.Sprintf("shows/%s/lists/%s/%s", *id, *t, *sort)
	} else {
		url = fmt.Sprintf("shows/%s/lists", *id)
	}

	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch lists url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.PersonalList{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found lists for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
