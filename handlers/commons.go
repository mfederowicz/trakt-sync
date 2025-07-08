// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// CommonInterface interface
type CommonInterface interface {
	FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, error)
	FetchShow(client *internal.Client, id *int) (*str.Show, error)
	FetchSeason(client *internal.Client, id *int) (*str.Season, error)
	FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error)
	FetchPerson(client *internal.Client, options *str.Options) (*str.Person, error)
	FetchList(client *internal.Client, options *str.Options) (*str.PersonalList, error)
	FetchComment(client *internal.Client, options *str.Options) (*str.Comment, error)
	FetchNotes(client *internal.Client, options *str.Options) (*str.Notes, error)
	FetchNotesItem(client *internal.Client, options *str.Options) (*str.NotesItem, error)
	FetchCommentItem(client *internal.Client, options *str.Options) (*str.CommentMediaItem, error)
	FetchCommentUserLikes(client *internal.Client, options *str.Options) (*str.CommentUserLike, error)
	FetchTrendingComments(client *internal.Client, options *str.Options) (*str.CommentItem, error)
	FetchRecentComments(client *internal.Client, options *str.Options) (*str.CommentItem, error)
	FetchUpdatedComments(client *internal.Client, options *str.Options) (*str.CommentItem, error)
	FetchMovieRecommendations(client *internal.Client, options *str.Options) ([]*str.Recommendation, error)
	FetchShowRecommendations(client *internal.Client, options *str.Options) ([]*str.Recommendation, error)
	UpdateComment(client *internal.Client, options *str.Options) (*str.Comment, error)
	UpdateNotes(client *internal.Client, options *str.Options) (*str.Notes, error)
	DeleteComment(client *internal.Client, options *str.Options) (*str.Comment, *str.Response, error)
	DeleteNotes(client *internal.Client, options *str.Options) (*str.Notes, *str.Response, error)
	HideMovieRecommendation(client *internal.Client, options *str.Options) (*str.Response, error)
	HideShowRecommendation(client *internal.Client, options *str.Options) (*str.Response, error)
	FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error)
	CheckSeasonNumber(code *string) (*string, *string, error)
	StartScrobble(client *internal.Client, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error)
	StopScrobble(client *internal.Client, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error)
	PauseScrobble(client *internal.Client, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error)
	CreateScrobble(client *internal.Client, options *str.Options) (*str.Scrobble, error)
	CreateScrobbleShowEpisode(client *internal.Client, options *str.Options) (*str.Scrobble, error)
	Checkin(client *internal.Client, checkin *str.Checkin) (*str.Checkin, *str.Response, error)
	CreateCheckin(client *internal.Client, options *str.Options) (*str.Checkin, error)
	CreateCheckinShowEpisode(client *internal.Client, options *str.Options) (*str.Checkin, error)
	Comment(client *internal.Client, comment *str.Comment) (*str.Comment, *str.Response, error)
	Notes(client *internal.Client, notes *str.Notes) (*str.Notes, *str.Response, error)
	Reply(client *internal.Client, id *int, comment *str.Comment) (*str.Comment, *str.Response, error)
	CheckSortAndTypes(options *str.Options) error
	CheckTypes(options *str.Options) error
	ToTimestamp(at string) (*str.Timestamp, error)
	ConvertDateString(date string, out string) string
	CurrentDateString(tz string) string
	DateLastDays(days int, tz string, full bool) string
	CheckDates(from string, to string, tz string) string
	ReadInput(items string) (*str.ItemsList, error)
	ConvertBytesToColletionItems(data []byte) (*str.ItemsList, error)
	UpdateHistoryListWithType(data []*str.ExportlistItem, strtype string) (*str.ItemsList, error)
	FetchHistoryListSeasons(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error)
}

// CommonLogic struct for common methods
type CommonLogic struct{}

// Media interface for helpers
type Media interface {
	str.Movie | str.Show | str.Episode | str.Season
}

// OnlySeasonsIDs helper for list of season ids object
func (c CommonLogic) OnlySeasonsIDs(items *[]str.ExportlistItem) *[]str.Season {
	result := onlyIDs[str.Season](*items)
	return &result
}

func (c CommonLogic) OnlyMoviesIDs(items *[]str.ExportlistItem) *[]str.Movie {
	result := onlyIDs[str.Movie](*items)
	return &result
}

func (c CommonLogic) OnlyShowsIDs(items *[]str.ExportlistItem) *[]str.Show {
	result := onlyIDs[str.Show](*items)
	return &result
}

func (c CommonLogic) OnlyEpisodesIDs(items *[]str.ExportlistItem) *[]str.Episode {
	result := onlyIDs[str.Episode](*items)
	return &result
}

// onlyIDs is a helper function to extract each type objects with only ids
func onlyIDs[T Media](items []str.ExportlistItem) []T {
	result := make([]T, 0, len(items))

	var zero T

	switch any(zero).(type) {
	case str.Movie:
		for _, item := range items {
			result = append(result, any(str.Movie{IDs: item.IDs}).(T))
		}
	case str.Show:
		for _, item := range items {
			if len(*item.Seasons) > 0 {
				updatedSeasons := SeasonsWithEpisodeNumbersOnly(item.Seasons)
				result = append(result, any(str.Show{IDs: item.IDs, Seasons: updatedSeasons}).(T))
			} else {
				result = append(result, any(str.Show{IDs: item.IDs}).(T))
			}
		}
	case str.Episode:
		for _, item := range items {
			result = append(result, any(str.Episode{IDs: item.IDs}).(T))
		}
	case str.Season:
		for _, item := range items {
			result = append(result, any(str.Season{IDs: item.IDs}).(T))
		}
	default:
		panic("unsupported type")
	}

	return result
}

// CreateItemsToRemove helper to create list of items to remove from history
func (c CommonLogic) CreateItemsToRemove(items *str.ItemsList) str.ItemsToRemove {

	return str.ItemsToRemove{
		Movies:   c.OnlyMoviesIDs(items.Movies),
		Shows:    c.OnlyShowsIDs(items.Shows),
		Seasons:  c.OnlySeasonsIDs(items.Seasons),
		Episodes: c.OnlyEpisodesIDs(items.Episodes),
	}

}

// CreateItemsToAdd helper to create list of items to add to history
func (c CommonLogic) CreateItemsToAdd(items *str.ItemsList) str.HistoryItems {
	movies := []str.Movie{}
	for _, m := range *items.Movies {
		movie := str.Movie{
			Title:     m.Title,
			Year:      m.Year,
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
		}
		movies = append(movies, movie)
	}
	shows := []str.Show{}
	for _, m := range *items.Shows {
		show := str.Show{
			Title:   m.Title,
			Year:    m.Year,
			IDs:     m.IDs,
			Seasons: m.Seasons,
		}
		shows = append(shows, show)
	}
	seasons := []str.Season{}
	for _, m := range *items.Seasons {
		season := str.Season{
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
		}
		seasons = append(seasons, season)
	}
	episodes := []str.Episode{}
	for _, m := range *items.Episodes {
		episode := str.Episode{
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
		}
		episodes = append(episodes, episode)
	}

	return str.HistoryItems{
		Movies:   &movies,
		Shows:    &shows,
		Seasons:  &seasons,
		Episodes: &episodes,
	}
}

// CreateCheckin helper function to create checkin object
func (c CommonLogic) CreateCheckin(client *internal.Client, options *str.Options) (*str.Checkin, error) {
	checkin := new(str.Checkin)
	connections, err := c.FetchUserConnections(client, options)
	if err != nil {
		return nil, fmt.Errorf(consts.UserConnectionsError, err)
	}
	checkin.Sharing = new(str.Sharing)
	checkin.Sharing.Tumblr = connections.Tumblr
	checkin.Sharing.Twitter = connections.Twitter
	checkin.Sharing.Mastodon = connections.Mastodon

	switch options.Action {
	case consts.Movie:
		movie, _, _ := c.FetchMovie(client, options)
		checkin.Movie = movie
	case consts.Episode:
		episode, _ := c.FetchEpisode(client, options)
		checkin.Episode = new(str.Episode)
		checkin.Episode.IDs = new(str.IDs)
		checkin.Episode.IDs.Trakt = episode.IDs.Trakt
	case consts.ShowEpisode:
		che, err := c.CreateCheckinShowEpisode(client, options)
		if err != nil {
			return nil, fmt.Errorf(consts.ShowEpisodeErr, err)
		}
		checkin.Show = che.Show
		checkin.Episode = che.Episode
	default:
		return nil, errors.New(consts.UnknownCheckinAction)
	}

	if len(options.Msg) > consts.ZeroValue {
		checkin.Message = &options.Msg
	}
	return checkin, nil
}

// CreateCheckinShowEpisode helper function to create checkin object for show episode
func (c CommonLogic) CreateCheckinShowEpisode(client *internal.Client, options *str.Options) (*str.Checkin, error) {
	checkin := new(str.Checkin)
	show, err := c.FetchShow(client, options)
	if err != nil {
		return nil, fmt.Errorf(consts.ShowErr, err)
	}
	checkin.Show = new(str.Show)
	checkin.Show = show
	checkin.Episode = new(str.Episode)

	if len(options.EpisodeCode) > consts.ZeroValue {
		season, number, err := c.CheckSeasonNumber(options.EpisodeCode)
		if err != nil {
			return nil, fmt.Errorf(consts.EpisodeCodeErr, err)
		}
		checkin.Episode.Season = season
		checkin.Episode.Number = number
	}
	if options.EpisodeAbs > consts.ZeroValue {
		checkin.Episode.NumberAbs = &options.EpisodeAbs
	}

	return checkin, nil
}

// CreateScrobble helper function to create scrobble object
func (c CommonLogic) CreateScrobble(client *internal.Client, options *str.Options) (*str.Scrobble, error) {
	scrobble := new(str.Scrobble)
	switch options.Type {
	case consts.Movie:
		movie, _, _ := c.FetchMovie(client, options)
		scrobble.Movie = movie
	case consts.Episode:
		episode, _ := c.FetchEpisode(client, options)
		scrobble.Episode = new(str.Episode)
		scrobble.Episode.IDs = new(str.IDs)
		scrobble.Episode.IDs.Trakt = episode.IDs.Trakt
	case consts.ShowEpisode:
		sc, err := c.CreateScrobbleShowEpisode(client, options)
		if err != nil {
			return nil, fmt.Errorf(consts.ScrobbleError, err)
		}
		scrobble.Episode = sc.Episode
		scrobble.Show = sc.Show
	}

	if options.Progress > consts.ZeroValue {
		scrobble.Progress = &options.Progress
	}

	return scrobble, nil
}

// CreateScrobbleShowEpisode helper function to create scrobble object
func (c CommonLogic) CreateScrobbleShowEpisode(client *internal.Client, options *str.Options) (*str.Scrobble, error) {
	scrobble := new(str.Scrobble)
	show, err := c.FetchShow(client, options)
	if err != nil {
		return nil, fmt.Errorf(consts.ShowErr, err)
	}
	scrobble.Show = new(str.Show)
	scrobble.Show = show
	scrobble.Episode = new(str.Episode)

	if len(options.EpisodeCode) > consts.ZeroValue {
		season, number, err := c.CheckSeasonNumber(options.EpisodeCode)
		if err != nil {
			return nil, fmt.Errorf(consts.EpisodeCodeErr, err)
		}
		scrobble.Episode.Season = season
		scrobble.Episode.Number = number
	}
	if options.EpisodeAbs > consts.ZeroValue {
		scrobble.Episode.NumberAbs = &options.EpisodeAbs
	}

	return scrobble, nil
}

// FetchMovie helper function to fetch movie object
func (*CommonLogic) FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, *str.Response, error) {
	movieID := options.InternalID
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Movies.GetMovie(
		client.BuildCtxFromOptions(options),
		&movieID,
		&opts,
	)

	return result, resp, err
}

// FetchShow helper function to fetch show object
func (*CommonLogic) FetchShow(client *internal.Client, options *str.Options) (*str.Show, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	showID := options.InternalID

	result, _, err := client.Shows.GetShow(
		client.BuildCtxFromOptions(options),
		&showID,
		&opts,
	)

	return result, err
}

// FetchSeason helper function to fetch season object
func (*CommonLogic) FetchSeason(client *internal.Client, options *str.Options) (*str.Season, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	seasonID := options.InternalID
	result, _, err := client.Seasons.GetSeason(
		client.BuildCtxFromOptions(options),
		&seasonID,
		&opts,
	)

	return result, err
}

// FetchEpisode helper function to fetch episode object
func (*CommonLogic) FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error) {
	episodeID := options.InternalID
	result, _, err := client.Episodes.GetEpisode(
		client.BuildCtxFromOptions(options),
		&episodeID,
	)

	return result, err
}

// FetchPerson helper function to fetch person object
func (*CommonLogic) FetchPerson(client *internal.Client, options *str.Options) (*str.Person, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	personID := options.InternalID
	result, _, err := client.People.GetSinglePerson(
		client.BuildCtxFromOptions(options),
		&personID,
		&opts,
	)

	return result, err
}

// FetchList helper function to fetch list object
func (*CommonLogic) FetchList(client *internal.Client, options *str.Options) (*str.PersonalList, error) {
	listID := options.InternalID
	result, _, err := client.Lists.GetList(
		client.BuildCtxFromOptions(options),
		&listID,
	)

	return result, err
}

// FetchComment helper function to fetch comment object
func (*CommonLogic) FetchComment(client *internal.Client, options *str.Options) (*str.Comment, error) {
	commentID := options.CommentID
	result, _, err := client.Comments.GetComment(
		client.BuildCtxFromOptions(options),
		&commentID,
	)

	return result, err
}

// FetchNotes helper function to fetch notes object
func (*CommonLogic) FetchNotes(client *internal.Client, options *str.Options) (*str.Notes, error) {
	notesID := options.InternalID
	result, _, err := client.Notes.GetNotes(
		client.BuildCtxFromOptions(options),
		&notesID,
	)

	return result, err
}

// FetchNotesItem helper function to fetch notes attached item object
func (*CommonLogic) FetchNotesItem(client *internal.Client, options *str.Options) (*str.NotesItem, error) {
	notesID := options.InternalID
	result, _, err := client.Notes.GetNotesItem(
		client.BuildCtxFromOptions(options),
		&notesID,
	)

	return result, err
}

// FetchCommentItem helper function to fetch comment media item object
func (*CommonLogic) FetchCommentItem(client *internal.Client, options *str.Options) (*str.CommentMediaItem, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	commentID := options.CommentID
	result, _, err := client.Comments.GetCommentItem(
		client.BuildCtxFromOptions(options),
		&commentID,
		&opts,
	)

	return result, err
}

// FetchCommentUserLikes helper function to fetch comment user like object
func (c *CommonLogic) FetchCommentUserLikes(client *internal.Client, options *str.Options, page int) ([]*str.CommentUserLike, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	commentID := options.CommentID
	list, resp, err := client.Comments.GetCommentUserLikes(
		client.BuildCtxFromOptions(options),
		&commentID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchCommentUserLikes(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}

// FetchTrendingComments helper function to fetch tending comments object
func (c *CommonLogic) FetchTrendingComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, IncludeReplies: options.IncludeReplies}
	commentType := options.CommentType
	strType := options.Type
	list, resp, err := client.Comments.GetTrendingComments(
		client.BuildCtxFromOptions(options),
		&commentType,
		&strType,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchTrendingComments(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}

// FetchRecentComments helper function to fetch recent comments object
func (c *CommonLogic) FetchRecentComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, IncludeReplies: options.IncludeReplies}
	commentType := options.CommentType
	strType := options.Type
	list, resp, err := client.Comments.GetRecentComments(
		client.BuildCtxFromOptions(options),
		&commentType,
		&strType,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchRecentComments(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}

// FetchUpdatedComments helper function to fetch updated comments object
func (c *CommonLogic) FetchUpdatedComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, IncludeReplies: options.IncludeReplies}
	commentType := options.CommentType
	strType := options.Type
	list, resp, err := client.Comments.GetUpdatedComments(
		client.BuildCtxFromOptions(options),
		&commentType,
		&strType,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchUpdatedComments(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}

// FetchMovieRecommendations helper function to fetch movie recommendations
func (c *CommonLogic) FetchMovieRecommendations(client *internal.Client, options *str.Options, page int) ([]*str.Recommendation, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, IgnoreCollected: options.IgnoreCollected, IgnoreWatchlisted: options.IgnoreWatchlisted}
	list, resp, err := client.Recommendations.GetMovieRecommendations(
		client.BuildCtxFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchMovieRecommendations(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}

// FetchShowRecommendations helper function to fetch movie recommendations
func (c *CommonLogic) FetchShowRecommendations(client *internal.Client, options *str.Options, page int) ([]*str.Recommendation, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, IgnoreCollected: options.IgnoreCollected, IgnoreWatchlisted: options.IgnoreWatchlisted}
	list, resp, err := client.Recommendations.GetShowRecommendations(
		client.BuildCtxFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchShowRecommendations(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}

// UpdateComment helper function to put comment object
func (*CommonLogic) UpdateComment(client *internal.Client, options *str.Options, comment *str.Comment) (*str.Comment, *str.Response, error) {
	commentID := options.CommentID
	result, resp, err := client.Comments.UpdateComment(
		client.BuildCtxFromOptions(options),
		&commentID,
		comment,
	)

	return result, resp, err
}

// UpdateNotes helper function to put notes object
func (*CommonLogic) UpdateNotes(client *internal.Client, options *str.Options, notes *str.Notes) (*str.Notes, *str.Response, error) {
	notesID := options.InternalID
	result, resp, err := client.Notes.UpdateNotes(
		client.BuildCtxFromOptions(options),
		&notesID,
		notes,
	)

	return result, resp, err
}

// DeleteComment helper function to delete comment object
func (*CommonLogic) DeleteComment(client *internal.Client, options *str.Options) (*str.Response, error) {
	commentID := options.CommentID
	resp, err := client.Comments.DeleteComment(
		client.BuildCtxFromOptions(options),
		&commentID,
	)

	return resp, err
}

// DeleteNotes helper function to delete notes object
func (*CommonLogic) DeleteNotes(client *internal.Client, options *str.Options) (*str.Response, error) {
	notesID := options.InternalID
	resp, err := client.Notes.DeleteNotes(
		client.BuildCtxFromOptions(options),
		&notesID,
	)

	return resp, err
}

// HideMovieRecommendation helper function to hide movie recommendations
func (*CommonLogic) HideMovieRecommendation(client *internal.Client, options *str.Options) (*str.Response, error) {
	movieID := options.InternalID
	resp, err := client.Recommendations.HideMovieRecommendation(
		client.BuildCtxFromOptions(options),
		&movieID,
	)

	return resp, err
}

// HideShowRecommendation helper function to hide show recommendations
func (*CommonLogic) HideShowRecommendation(client *internal.Client, options *str.Options) (*str.Response, error) {
	showID := options.InternalID
	resp, err := client.Recommendations.HideShowRecommendation(
		client.BuildCtxFromOptions(options),
		&showID,
	)

	return resp, err
}

// FetchUserConnections helper function to fetch connections object
func (*CommonLogic) FetchUserConnections(client *internal.Client, options *str.Options) (*str.Connections, error) {
	result, _, err := client.Users.RetrieveSettings(
		client.BuildCtxFromOptions(options),
	)
	if err != nil {
		return nil, fmt.Errorf(consts.UserSettingsError, err)
	}

	return result.Connections, err
}

// StartScrobble helper function to start scrobble
func (*CommonLogic) StartScrobble(client *internal.Client, scrobble *str.Scrobble, options *str.Options) (*str.Scrobble, *str.Response, error) {
	result, resp, err := client.Scrobble.StartScrobble(
		client.BuildCtxFromOptions(options),
		scrobble,
	)

	return result, resp, err
}

// StopScrobble helper function to stop scrobble
func (*CommonLogic) StopScrobble(client *internal.Client, scrobble *str.Scrobble, options *str.Options) (*str.Scrobble, *str.Response, error) {
	result, resp, err := client.Scrobble.StopScrobble(
		client.BuildCtxFromOptions(options),
		scrobble,
	)

	return result, resp, err
}

// PauseScrobble helper function to pause scrobble
func (*CommonLogic) PauseScrobble(client *internal.Client, scrobble *str.Scrobble, options *str.Options) (*str.Scrobble, *str.Response, error) {
	result, resp, err := client.Scrobble.PauseScrobble(
		client.BuildCtxFromOptions(options),
		scrobble,
	)

	return result, resp, err
}

// Checkin helper function to post checkin object
func (*CommonLogic) Checkin(client *internal.Client, checkin *str.Checkin, options *str.Options) (*str.Checkin, *str.Response, error) {
	result, resp, err := client.Checkin.CheckintoAnItem(
		client.BuildCtxFromOptions(options),
		checkin,
	)

	return result, resp, err
}

// Comment helper function to post comment object
func (*CommonLogic) Comment(client *internal.Client, comment *str.Comment, options *str.Options) (*str.Comment, *str.Response, error) {
	result, resp, err := client.Comments.PostAComment(
		client.BuildCtxFromOptions(options),
		comment,
	)
	return result, resp, err
}

// Notes helper function to post notes object
func (*CommonLogic) Notes(client *internal.Client, notes *str.Notes, options *str.Options) (*str.Notes, *str.Response, error) {
	result, resp, err := client.Notes.AddNotes(
		client.BuildCtxFromOptions(options),
		notes,
	)
	return result, resp, err
}

// Reply helper function to post reply object
func (*CommonLogic) Reply(client *internal.Client, id *int, reply *str.Comment, options *str.Options) (*str.Comment, *str.Response, error) {
	result, resp, err := client.Comments.ReplyAComment(
		client.BuildCtxFromOptions(options),
		id,
		reply,
	)
	return result, resp, err
}

// CheckSeasonNumber helper function to convert string to season and episode
func (*CommonLogic) CheckSeasonNumber(code string) (season *int, episode *int, err error) {
	if len(code) < int(consts.MinSeasonNumberLength) {
		return nil, nil, errors.New("invalid length")
	}

	if parts := strings.Split(code, "x"); len(parts) == consts.TwoValue {
		season, _ := strconv.Atoi(parts[0])
		episode, _ := strconv.Atoi(parts[1])
		return &season, &episode, nil
	}

	return nil, nil, errors.New("invalid format")
}

// CheckSortAndTypes helper function to validate sort and type field depends on module
func (*CommonLogic) CheckSortAndTypes(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	_, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	prefix := options.Module + ":" + options.Action
	if !cfg.IsValidConfigType(cfg.ModuleActionConfig[prefix].Type, options.Type) {
		return fmt.Errorf("not found type for module '%s'", options.Module)
	}

	if !cfg.IsValidConfigType(cfg.ModuleActionConfig[prefix].Sort, options.Sort) {
		return fmt.Errorf("not found sort for module '%s'", options.Module)
	}

	// Check id_type values
	return nil
}

// CheckTypes helper function to validate type field depends on module
func (CommonLogic) CheckTypes(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	_, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	prefix := options.Module + ":" + options.Action
	if !cfg.IsValidConfigType(cfg.ModuleActionConfig[prefix].Type, options.Type) {
		return fmt.Errorf("not found type for module '%s'", options.Module)
	}

	// Check id_type values
	return nil
}

// ValidPrivacy helper function to validate privacy field depends on module
func (*CommonLogic) ValidPrivacy(options *str.Options) error {
	// Check if the provided module exists in ModuleConfig
	_, ok := cfg.ModuleConfig[options.Module]
	if !ok {
		return fmt.Errorf("not found config for module '%s'", options.Module)
	}
	prefix := options.Module + ":" + options.Action
	if len(options.Privacy) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.ModuleActionConfig[prefix].Privacy, options.Privacy) {
		return fmt.Errorf("invalid privacy '%s' for module '%s'", options.Privacy, options.Module)
	}
	return nil
}

// GenActionsUsage prints a usage message when an invalid action is provided.
func (*CommonLogic) GenActionsUsage(name string, actions []string) {
	printer.Println("Usage: ./trakt-sync " + name + " -a [action]")
	printer.Println("Available actions:")
	for _, action := range actions {
		printer.Printf(consts.ListItem, action)
	}
}

// GenTypeUsage prints a usage message when an invalid type is provided.
func (*CommonLogic) GenTypeUsage(name string, types []string) {
	printer.Println("Usage: ./trakt-sync " + name + " -t [type]")
	printer.Println("Available types:")
	for _, t := range types {
		printer.Printf(consts.ListItem, t)
	}
}

// GenActionTypeUsage prints a usage message when an invalid type for action is provided.
func (*CommonLogic) GenActionTypeUsage(options *str.Options, types []string) {
	printer.Println("Usage: ./trakt-sync " + options.Module + " -a " + options.Action + " -t [type]")
	printer.Println("Available types:")
	for _, t := range types {
		printer.Printf(consts.ListItem, t)
	}
}

// GenActionTypeItemUsage prints a usage message when an invalid item for type is provided.
func (*CommonLogic) GenActionTypeItemUsage(options *str.Options, items []string) {
	printer.Println("Usage: ./trakt-sync " + options.Module + " -a " + options.Action + " -t " + options.Type + " -item [item]")
	printer.Println("Available items:")
	for _, t := range items {
		printer.Printf("  - %s\n", t)
	}
}

// GetHandlerForMap choose handler from map
func (*CommonLogic) GetHandlerForMap(action string, allHandlers map[string]Handler) (Handler, error) {
	// Lookup and execute handler
	if handler, found := allHandlers[action]; found {
		return handler, nil
	}

	return nil, errors.New("unknown handler")
}

// ConvertDateString takes a date string and converts it to date time format,
// if empty return current date
func (CommonLogic) ConvertDateString(dateStr string, outputFormat string, tz string, full bool) string {
	// Parse the input date string using YYYY-MM-DD
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		parsedDate = time.Now()
	}

	// Get the current time
	currentTime := time.Now().UTC()

	// Combine the parsed date with the current time's hour, minute, second
	finalDateTime := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
		currentTime.Nanosecond(),
		currentTime.Location(),
	)

	if tz != time.UTC.String() {
		loc, _ := time.LoadLocation(tz)
		finalDateTime = finalDateTime.In(loc)
	}

	if full {
		finalDateTime = finalDateTime.Truncate(time.Hour)
	}
	// Format the parsed time into the output format
	formattedDate := finalDateTime.Format(outputFormat)

	return formattedDate
}

// ToTimestamp convert date time string to Timestamp
func (CommonLogic) ToTimestamp(at string) *str.Timestamp {
	// Parse the input date string using YYYY-MM-DD
	parsedDate, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return &str.Timestamp{}
	}

	return &str.Timestamp{Time: parsedDate.UTC()}
}

// CurrentDateString return current date string from user timezone
func (CommonLogic) CurrentDateString(tz string, full bool) string {
	// Get the current time
	currentTime := time.Now().UTC()

	if tz != time.UTC.String() {
		loc, _ := time.LoadLocation(tz)
		currentTime = currentTime.In(loc)
	}

	if full {
		currentTime = currentTime.Truncate(time.Hour)
	}

	return currentTime.Format(time.RFC3339)
}

// DateLastDays return last x days from user timezone
func (CommonLogic) DateLastDays(days int, tz string, full bool) string {
	// Get the current time
	currentTime := time.Now().UTC()

	if tz != time.UTC.String() {
		loc, _ := time.LoadLocation(tz)
		currentTime = currentTime.In(loc)
	}

	past := currentTime.AddDate(0, 0, -days)
	if full {
		past = past.Truncate(time.Hour)
	}

	return past.Format(time.RFC3339)
}

// CheckDates if dates are ok
func (CommonLogic) CheckDates(from string, to string, tz string) error {
	// Load location
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	// Get current full hour
	now := time.Now().In(loc)
	fullHour := now.Truncate(time.Hour)

	// Parse from and to if they are present
	var fromTime, toTime time.Time
	if from != "" {
		fromTime, err = time.ParseInLocation(consts.DefaultStartDateFormat, from, loc)
		if err != nil {
			return fmt.Errorf("invalid from date: %w", err)
		}
		if fromTime.After(fullHour) {
			return fmt.Errorf("'from' date must not be later than the current full hour")
		}
	}

	if to != "" {
		toTime, err = time.ParseInLocation(consts.DefaultStartDateFormat, to, loc)
		if err != nil {
			return fmt.Errorf("invalid to date: %w", err)
		}
		if toTime.After(fullHour) {
			return fmt.Errorf("'to' date must not be later than:%s", fullHour.Format(consts.DefaultStartDateFormat))
		}
	}

	// If both are present, validate range
	if from != "" && to != "" {
		if fromTime.After(toTime) {
			return fmt.Errorf("'from' date must be earlier than or equal to 'to' date")
		}
	}

	return nil
}

// ConvertBytesToItemsList convert bytes to struct
func (c *CommonLogic) ConvertBytesToItemsList(data []byte, action string, stype string) (*str.ItemsList, error) {
	var list []*str.ExportlistItem
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	items := new(str.ItemsList)
	items.Movies = &[]str.ExportlistItem{}
	items.Shows = &[]str.ExportlistItem{}
	items.Seasons = &[]str.ExportlistItem{}
	items.Episodes = &[]str.ExportlistItem{}
	items.IDs = &[]int64{}
	switch action {
	case consts.AddToHistory:
		items = c.ListToHistoryItems(items, list, stype)
	case consts.AddToCollection:
		items = c.ListToCollectionItems(items, list, stype)
	default:
		return nil, errors.New(consts.UnknownItemsListType)
	}

	return items.Uniq(), nil
}

func (c *CommonLogic) ListToCollectionItems(items *str.ItemsList, list []*str.ExportlistItem, stype string) *str.ItemsList {
	for _, val := range list {
		if val.ID != nil {
			*items.IDs = append(*items.IDs, *val.ID)
		}

		if val.Movie != nil {
			e := str.ExportlistItem{}
			e.Title = val.Movie.Title
			e.Year = val.Movie.Year
			e.IDs = val.Movie.IDs
			e.UpdateCollectedData(val)
			*items.Movies = append(*items.Movies, e)
		}
		if val.Show != nil {
			e := str.ExportlistItem{}
			e.Title = val.Show.Title
			e.Year = val.Show.Year
			e.IDs = val.Show.IDs
			e.UpdateCollectedData(val)
			*items.Shows = append(*items.Shows, e)
		}
		if val.Season != nil {
			e := str.ExportlistItem{}
			e.IDs = val.Season.IDs
			val.Season.UpdateCollectedData(val)
			*items.Seasons = append(*items.Seasons, e)

		}
		if val.Episode != nil {
			e := str.ExportlistItem{}
			e.IDs = val.Episode.IDs
			e.UpdateCollectedData(val)
			*items.Episodes = append(*items.Episodes, e)

		}

	}

	return items

}

// ListToHistoryItems convert list to history items to struct
func (c *CommonLogic) ListToHistoryItems(items *str.ItemsList, list []*str.ExportlistItem, stype string) *str.ItemsList {
	// Group by movie
	moviesMap := make(map[int64]str.OutputMovie)
	// Group by show
	showsMap := make(map[int64]str.OutputShow)
	// Group by season
	seasonsMap := make(map[int64]str.OutputSeason)
	// Group by episodes
	episodesMap := make(map[int64]str.OutputEpisode)

	for _, item := range list {
		switch stype {
		case consts.Movies:
			movieID := item.Movie.IDs.Trakt
			if item.ID != nil {
				*items.IDs = append(*items.IDs, *item.ID)
			}
			if movieID != nil {
				movie, exists := moviesMap[*movieID]
				if !exists {
					movie = str.OutputMovie{
						Title:     item.Movie.Title,
						Year:      item.Movie.Year,
						IDs:       item.Movie.IDs,
						WatchedAt: item.WatchedAt,
					}
				}

				moviesMap[*movieID] = movie
			}

		case consts.Shows:

			showID := item.Show.IDs.Trakt
			if showID != nil {
				showIDnp := *showID
				if item.ID != nil {
					*items.IDs = append(*items.IDs, *item.ID)
				}
				show, exists := showsMap[showIDnp]
				if !exists {
					show = str.OutputShow{
						Title:     item.Show.Title,
						Year:      item.Show.Year,
						IDs:       item.Show.IDs,
						Seasons:   &[]str.Season{},
						WatchedAt: item.WatchedAt,
					}
				}
				// Initialize seasons if nil
				if show.Seasons == nil {
					empty := []str.Season{}
					show.Seasons = &empty
				}

				// Find or create season
				var season *str.Season

				for i := range *show.Seasons {
					if *(*show.Seasons)[i].Number == *item.Episode.Season {
						season = &(*show.Seasons)[i]
						break
					}
				}

				if season == nil {
					newSeason := str.Season{
						Number: item.Episode.Season,
					}
					*show.Seasons = append(*show.Seasons, newSeason)
					season = &(*show.Seasons)[len(*show.Seasons)-1]
					season.Episodes = &[]str.Episode{}
				}

				// Add episode
				newEpisode := &str.Episode{Number: item.Episode.Number, WatchedAt: item.WatchedAt}
				*season.Episodes = append(*season.Episodes, *newEpisode)

				showsMap[showIDnp] = show
			}

		case consts.Seasons:
			seasonID := item.Season.IDs.Trakt
			if seasonID != nil {
				season, exists := seasonsMap[*seasonID]
				if !exists {
					season = str.OutputSeason{
						Title:     item.Season.Title,
						IDs:       item.Season.IDs,
						WatchedAt: item.WatchedAt,
					}
				}
				seasonsMap[*seasonID] = season
			}
		case consts.Episodes:
			episodeID := item.Episode.IDs.Trakt
			if episodeID != nil {
				episode, exists := episodesMap[*episodeID]
				if !exists {
					episode = str.OutputEpisode{
						Title:     item.Episode.Title,
						IDs:       item.Episode.IDs,
						WatchedAt: item.WatchedAt,
					}
				}
				episodesMap[*episodeID] = episode
			}
		default:
			return &str.ItemsList{}
		}
	}

	for _, s := range episodesMap {
		*items.Episodes = append(*items.Episodes, str.ExportlistItem{
			Episode:   &str.Episode{Title: s.Title, IDs: s.IDs},
			WatchedAt: s.WatchedAt,
			IDs:       s.IDs,
		})
	}
	for _, s := range seasonsMap {
		*items.Seasons = append(*items.Seasons, str.ExportlistItem{
			Season:    &str.Season{Title: s.Title, IDs: s.IDs},
			WatchedAt: s.WatchedAt,
			IDs:       s.IDs,
		})
	}
	for _, s := range moviesMap {
		*items.Movies = append(*items.Movies, str.ExportlistItem{
			Movie:     &str.Movie{Title: s.Title, Year: s.Year, IDs: s.IDs},
			WatchedAt: s.WatchedAt,
			IDs:       s.IDs,
		})
	}

	for _, s := range showsMap {
		*items.Shows = append(*items.Shows, str.ExportlistItem{
			Show:      &str.Show{Title: s.Title, Year: s.Year, IDs: s.IDs},
			Seasons:   s.Seasons,
			WatchedAt: s.WatchedAt,
			IDs:       s.IDs,
		})

	}

	return items
}

// ReadInput read data from stdin or from file
func (c CommonLogic) ReadInput(options str.Options) (*str.ItemsList, error) {
	filePath := options.Items
	if filePath != consts.EmptyString {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		return c.ConvertBytesToItemsList(data, options.Action, options.Type)
	}

	// Check if there's data in stdin to avoid blocking
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}

	// os.ModeCharDevice means no data is being piped (stdin is a terminal)
	if fi.Mode()&os.ModeCharDevice != 0 {
		return nil, fmt.Errorf("no --file provided and no data piped to stdin")
	}

	// Read all data from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}

	return c.ConvertBytesToItemsList(data, options.Action, options.Type)
}

// FetchHistoryList returns movies and episodes that a user has watched, sorted by most recent.
func (c CommonLogic) FetchHistoryList(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, StartAt: options.StartDate, EndAt: options.EndDate, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetWatchedHistory(
		client.BuildCtxFromOptions(options),
		&options.TraktID,
		&options.Type,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := c.FetchHistoryList(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

func (c CommonLogic) UpdateHistoryListWithType(data []*str.ExportlistItem, strtype *string) []*str.ExportlistItem {
	list := []*str.ExportlistItem{}

	newType := consts.EmptyString
	switch *strtype {
	case consts.Shows:
		newType = consts.Show
	case consts.Episodes:
		newType = consts.Episode
	case consts.Movies:
		newType = consts.Movie
	case consts.Seasons:
		newType = consts.Season
	default:
		newType = consts.Movie
	}

	for _, item := range data {
		item.Type = &newType
		list = append(list, item)
	}

	return list
}

func (c CommonLogic) FetchHistoryListSeasons(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {

	// fetch history shows
	options.Type = consts.Shows
	_, err := c.FetchHistoryList(client, options, page)

	if err != nil {
		return nil, err
	}
	// for _, val := range episodes {
	// 	//fmt.Println(*val.ID)
	// }

	// collected := []str.Season{}
	// opts := uri.ListOptions{Extended: options.ExtendedInfo}
	// for _, val := range shows {
	// 	time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
	//
	// 	seasonsNumbers := []int{}
	// 	for _, sitem := range *val.Seasons {
	// 		seasonsNumbers = append(seasonsNumbers, *sitem.Number)
	// 	}
	//
	// 	seasons, _, err := client.Shows.GetAllSeasonsForShow(client.BuildCtxFromOptions(options), val.Show.IDs.Slug, &opts)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, sitem := range seasons {
	// 		if slices.Contains(seasonsNumbers, *sitem.Number) {
	// 			s := str.Season{}
	// 			s.IDs = sitem.IDs
	// 			collected = append(collected, s)
	// 		}
	// 	}
	// }
	// fmt.Println(collected)

	return nil, nil

}

// Ptr is a helper routine that allocates a new T value
// to store v and returns a pointer to it.
func Ptr[T any](v T) *T {
	return &v
}

// SeasonsWithEpisodeNumbersOnly is a helper function to make seasons lists with episodes contains numbers
func SeasonsWithEpisodeNumbersOnly(src *[]str.Season) *[]str.Season {

	if src == nil {
		return nil
	}

	seasonsCopy := make([]str.Season, len(*src))

	for i, season := range *src {
		seasonsCopy[i] = season // shallow copy

		if season.Episodes != nil {
			eps := make([]str.Episode, len(*season.Episodes))

			for j, ep := range *season.Episodes {
				eps[j] = str.Episode{
					Number: ep.Number,
				}
			}

			seasonsCopy[i].Episodes = &eps
		}
	}

	return &seasonsCopy
}
