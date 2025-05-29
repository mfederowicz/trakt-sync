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

// GetShowCollectionProgress Returns collection progress for a show including details on all aired seasons and episodes.
// API docs: https://trakt.docs.apiary.io/#reference/shows/collection-progress/get-show-collection-progress
func (s *ShowsService) GetShowCollectionProgress(ctx context.Context, id *string, opts *uri.ListOptions) (*str.CollectionProgress, error) {
	var url string
	url = fmt.Sprintf("shows/%s/progress/collection", *id)

	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, err
	}

	printer.Println("fetch collection progress url:" + url)

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	collection := new(str.CollectionProgress)
	_, err = s.client.Do(ctx, req, &collection)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

// GetShowWatchedProgress Returns watched progress for a show including details on all aired seasons and episodes.
// API docs: https://trakt.docs.apiary.io/#reference/shows/watched-progress/get-show-watched-progress
func (s *ShowsService) GetShowWatchedProgress(ctx context.Context, id *string, opts *uri.ListOptions) (*str.WatchedProgress, error) {
	var url string
	url = fmt.Sprintf("shows/%s/progress/watched", *id)

	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, err
	}

	printer.Println("fetch watched progress url:" + url)

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	watched := new(str.WatchedProgress)
	_, err = s.client.Do(ctx, req, &watched)

	if err != nil {
		return nil, err
	}

	return watched, nil
}

// ResetShowProgress Reset a show's progress when the user started re-watching the show.
// API docs:https://trakt.docs.apiary.io/#reference/shows/reset-watched-progress/reset-show-progress
func (s *ShowsService) ResetShowProgress(ctx context.Context, id *string, progress *str.WatchedProgress) (*str.WatchedProgress, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/progress/watched/reset", *id)
	req, err := s.client.NewRequest(http.MethodPost, url, progress)
	if err != nil {
		return nil, nil, err
	}

	respProgress := new(str.WatchedProgress)
	resp, err := s.client.Do(ctx, req, respProgress)
	if err != nil {
		return nil, nil, err
	}

	return respProgress, resp, nil
}

// UndoResetShowProgress Undo the reset and have watched progress use all watched history for the show.
// API docs:https://trakt.docs.apiary.io/#reference/shows/reset-watched-progress/undo-reset-show-progress
func (s *ShowsService) UndoResetShowProgress(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("shows/%s/progress/watched/reset", *id)
	printer.Println("undo reset watched progress")
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetShowRatings Returns rating (between 0 and 10) and distribution for a show.
// API docs: https://trakt.docs.apiary.io/#reference/shows/ratings/get-show-ratings
func (s *ShowsService) GetShowRatings(ctx context.Context, id *string) (*str.ShowRatings, *str.Response, error) {
	url := fmt.Sprintf("shows/%s/ratings", *id)
	printer.Println("fetch ratings url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.ShowRatings)
	resp, err := s.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found ratings for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch ratings err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetRelatedShows Returns related and similar shows.
// API docs: https://trakt.docs.apiary.io/#reference/shows/related/get-related-shows
func (s *ShowsService) GetRelatedShows(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.Show, *str.Response, error) {
	var url string
	url = fmt.Sprintf("shows/%s/related", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch related url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Show{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found related for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch related err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetShowStats Returns lots of show stats.
// API docs: https://trakt.docs.apiary.io/#reference/shows/stats/get-show-stats
func (s *ShowsService) GetShowStats(ctx context.Context, id *string) (*str.ShowStats, *str.Response, error) {
	url := fmt.Sprintf("shows/%s/stats", *id)
	printer.Println("fetch stats url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.ShowStats)
	resp, err := s.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found stats for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch stats err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetShowStudios Returns all studios for show.
// API docs: https://trakt.docs.apiary.io/#reference/shows/studios/get-show-studios
func (s *ShowsService) GetShowStudios(ctx context.Context, id *string) ([]*str.Studio, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/studios", *id)
	printer.Println("fetch studios url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Studio{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found studios for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch studios err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetShowWatching Returns all users watching this show right now.
// API docs:  https://trakt.docs.apiary.io/#reference/shows/studios/get-users-watching-right-now
func (s *ShowsService) GetShowWatching(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.UserProfile, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/watching", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch watching url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.UserProfile{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found watching for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch watching err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetShowVideos Returns all videos including trailers, teasers, clips, and featurettes.
// API docs: https://trakt.docs.apiary.io/#reference/shows/videos/get-all-videos
func (s *ShowsService) GetShowVideos(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.Video, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/videos", *id)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch video url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Video{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found video for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch video err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// RefreshShowMetadata Queue this show for a full metadata and image refresh.
// It might take up to 8 hours for the updated metadata to be availabe through the API.
// API docs: https://trakt.docs.apiary.io/#reference/shows/refresh/refresh-show-metadata
func (s *ShowsService) RefreshShowMetadata(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("shows/%s/refresh", *id)
	printer.Println("refresh show:" + url)
	req, err := s.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetNextEpisode Returns the next scheduled to air episode.
// API docs: https://trakt.docs.apiary.io/#reference/shows/next-episode/get-next-episode
func (s *ShowsService) GetNextEpisode(ctx context.Context, id *string, opts *uri.ListOptions) (*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/next_episode", *id)
	url, err := uri.AddQuery(url, opts)
	printer.Println("fetch next episode url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.Episode)
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch next episode err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetLastEpisode Returns the most recently aired episode.
// API docs: https://trakt.docs.apiary.io/#reference/shows/last-episode/get-last-episode
func (s *ShowsService) GetLastEpisode(ctx context.Context, id *string, opts *uri.ListOptions) (*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/last_episode", *id)
	url, err := uri.AddQuery(url, opts)
	printer.Println("fetch last episode url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.Episode)
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch last episode err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetAllSeasonsForShow Returns all seasons for a show including the number of episodes in each season.
// API docs: https://trakt.docs.apiary.io/#reference/seasons/summary/get-all-seasons-for-a-show
func (s *ShowsService) GetAllSeasonsForShow(ctx context.Context, id *string, opts *uri.ListOptions) ([]*str.Season, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons", *id)
	url, err := uri.AddQuery(url, opts)
	printer.Println("fetch all seasons url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := []*str.Season{}
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch seasons err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSingleSeasonsForShow Returns a single seasons for a show.
// API docs: https://trakt.docs.apiary.io/#reference/seasons/season/get-single-seasons-for-a-show
func (s *ShowsService) GetSingleSeasonsForShow(ctx context.Context, id *string, season *int, opts *uri.ListOptions) (*str.Season, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons/%d/info", *id, *season)
	url, err := uri.AddQuery(url, opts)
	printer.Println("fetch single seasons url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Season)
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch seasons err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetAllEpisodesForSingleSeason Returns a single seasons for a show.
// API docs: https://trakt.docs.apiary.io/#reference/seasons/episodes/get-all-episodes-for-a-single-season
func (s *ShowsService) GetAllEpisodesForSingleSeason(ctx context.Context, id *string, season *int, opts *uri.ListOptions) ([]*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons/%d", *id, *season)
	url, err := uri.AddQuery(url, opts)
	printer.Println("fetch season episodes url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := []*str.Episode{}
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch season episodes err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetAllSeasonTranslations Returns all translations for an season, including language and translated values for title and overview.
// API docs: https://trakt.docs.apiary.io/#reference/seasons/episodes/get-all-season-translations
func (s *ShowsService) GetAllSeasonTranslations(ctx context.Context, id *string, season *int, language *string, opts *uri.ListOptions) ([]*str.Translation, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons/%d/translations", *id, *season)

	if len(*language) > consts.ZeroValue {
		url = fmt.Sprintf("shows/%s/seasons/%d/translations/%s", *id, *season, *language)
	}

	url, err := uri.AddQuery(url, opts)
	printer.Println("fetch season translations url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := []*str.Translation{}
	resp, err := s.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch season translations err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetAllSeasonComments Returns all top level comments for a season.
// By default, the newest comments are returned first.
// Other sorting options include oldest, most likes, most replies, highest rated, lowest rated, and most plays..
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/comments/get-all-season-comments
func (s *ShowsService) GetAllSeasonComments(ctx context.Context, id *string, season *int, sort *string, opts *uri.ListOptions) ([]*str.Comment, *str.Response, error) {
	var url string
	if *sort != consts.EmptyString {
		url = fmt.Sprintf("shows/%s/seasons/%d/comments/%s", *id, *season, *sort)
	} else {
		url = fmt.Sprintf("shows/%s/seasons/%d/comments", *id, *season)
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

// GetListsContainingSeason Returns all lists that contain this season.
// By default, personal lists are returned sorted by the most popular.
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/lists/get-lists-containing-this-season
func (s *ShowsService) GetListsContainingSeason(ctx context.Context, id *string, season *int, t *string, sort *string, opts *uri.ListOptions) ([]*str.PersonalList, *str.Response, error) {
	var url string
	if *t != consts.EmptyString && *sort != consts.EmptyString {
		url = fmt.Sprintf("shows/%s/seasons/%d/lists/%s/%s", *id, *season, *t, *sort)
	} else {
		url = fmt.Sprintf("shows/%s/seasons/%d/lists", *id, *season)
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

// GetAllPeopleForSeason Returns all cast and crew for a season.
// Each cast member will have a characters array and a standard person object.
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/people/get-all-people-for-a-season
func (s *ShowsService) GetAllPeopleForSeason(ctx context.Context, id *string, season *int, opts *uri.ListOptions) (*str.SeasonPeople, *str.Response, error) {
	var url string

	url = fmt.Sprintf("shows/%s/seasons/%d/people", *id, *season)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch season people url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.SeasonPeople)
	resp, err := s.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found season people for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch season people err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSeasonRatings Returns rating (between 0 and 10) and distribution for a season.
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/ratings/get-season-ratings
func (s *ShowsService) GetSeasonRatings(ctx context.Context, id *string, season *int) (*str.SeasonRatings, *str.Response, error) {
	url := fmt.Sprintf("shows/%s/seasons/%d/ratings", *id, *season)
	printer.Println("fetch season ratings url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.SeasonRatings)
	resp, err := s.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found season ratings for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch seasons ratings err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSeasonStats Returns lots of season stats.
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/stats/get-season-stats
func (s *ShowsService) GetSeasonStats(ctx context.Context, id *string, season *int) (*str.SeasonStats, *str.Response, error) {
	url := fmt.Sprintf("shows/%s/seasons/%d/stats", *id, *season)
	printer.Println("fetch season stats url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	result := new(str.SeasonStats)
	resp, err := s.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found season stats for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch season stats err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSeasonsWatching Returns all users watching this season right now.
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/watching/get-users-watching-right-now
func (s *ShowsService) GetSeasonsWatching(ctx context.Context, id *string, season *int, opts *uri.ListOptions) ([]*str.UserProfile, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons/%d/watching", *id, *season)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch seasons watching url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.UserProfile{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found season watching for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch season watching err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetSeasonsVideos Returns all videos including trailers, teasers, clips, and featurettes.
//
// API docs: https://trakt.docs.apiary.io/#reference/seasons/videos/get-all-videos
func (s *ShowsService) GetSeasonsVideos(ctx context.Context, id *string, season *int, opts *uri.ListOptions) ([]*str.Video, *str.Response, error) {
	var url = fmt.Sprintf("shows/%s/seasons/%d/videos", *id, *season)
	url, err := uri.AddQuery(url, opts)

	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch season video url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Video{}
	resp, err := s.client.Do(ctx, req, &list)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("not found season video for id/slug:%s", *id)
	}

	if err != nil {
		printer.Println("fetch season video err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
