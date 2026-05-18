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
	ApproveFollowRequest(client *internal.Client, options *str.Options) (*str.FollowRequest, *str.Response, error)
	CheckDates(from string, to string, tz string) error
	CheckSeasonNumber(code string) (*int, *int, error)
	CheckSortAndTypes(options *str.Options) error
	CheckTypes(options *str.Options) error
	Checkin(client *internal.Client, checkin *str.Checkin, options *str.Options) (*str.Checkin, *str.Response, error)
	Comment(client *internal.Client, comment *str.Comment, options *str.Options) (*str.Comment, *str.Response, error)
	ConvertDateString(dateStr string, outputFormat string, tz string, full bool) string
	CreateCheckin(client *internal.Client, options *str.Options) (*str.Checkin, error)
	CreateCheckinShowEpisode(client *internal.Client, options *str.Options) (*str.Checkin, error)
	CreateItemsToAdd(items *str.ItemsList) str.HistoryItems
	CreateItemsToAddRatings(items *str.ItemsList) str.RatingItems
	CreateItemsToHidden(section string, items *str.ItemsList) str.HistoryItems
	CreateItemsToRemove(items *str.ItemsList) str.ItemsToRemove
	CreateItemsToReorder(items *str.ItemsList) str.ItemsToReorder
	CreateScrobble(client *internal.Client, options *str.Options) (*str.Scrobble, error)
	CreateScrobbleShowEpisode(client *internal.Client, options *str.Options) (*str.Scrobble, error)
	CurrentDateString(tz string, full bool) string
	DateLastDays(days int, tz string, full bool) string
	DeleteComment(client *internal.Client, options *str.Options) (*str.Response, error)
	DeleteNotes(client *internal.Client, options *str.Options) (*str.Response, error)
	DenyFollowRequest(client *internal.Client, options *str.Options) (*str.FollowRequest, *str.Response, error)
	FetchComment(client *internal.Client, options *str.Options) (*str.Comment, error)
	FetchCommentItem(client *internal.Client, options *str.Options) (*str.CommentMediaItem, error)
	FetchCommentUserLikes(client *internal.Client, options *str.Options, page int) ([]*str.CommentUserLike, error)
	FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error)
	FetchFavorites(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error)
	FetchFollowRequests(client *internal.Client, options *str.Options) ([]*str.FollowRequest, error)
	FetchHistoryList(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error)
	FetchList(client *internal.Client, options *str.Options) (*str.PersonalList, error)
	FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, *str.Response, error)
	FetchMovieRecommendations(client *internal.Client, options *str.Options, page int) ([]*str.Recommendation, error)
	FetchNotes(client *internal.Client, options *str.Options) (*str.Notes, error)
	FetchNotesItem(client *internal.Client, options *str.Options) (*str.NotesItem, error)
	FetchPendingFollowingRequests(client *internal.Client, options *str.Options) ([]*str.FollowRequest, error)
	FetchPerson(client *internal.Client, options *str.Options) (*str.Person, error)
	FetchRatings(client *internal.Client, options *str.Options, page int) ([]*str.RatingListItem, error)
	FetchRecentComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error)
	FetchSeason(client *internal.Client, options *str.Options) (*str.Season, error)
	FetchShow(client *internal.Client, options *str.Options) (*str.Show, error)
	FetchShowRecommendations(client *internal.Client, options *str.Options, page int) ([]*str.Recommendation, error)
	FetchTrendingComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error)
	FetchUpdatedComments(client *internal.Client, options *str.Options, page int) ([]*str.CommentItem, error)
	FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error)
	FetchUsersHiddenItems(client *internal.Client, options *str.Options, page int) ([]*str.HiddenItem, error)
	FetchWatchlist(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error)
	GenActionTypeItemUsage(options *str.Options, items []string)
	GenActionTypeUsage(options *str.Options, types []string)
	GenActionsUsage(name string, actions []string)
	GenTypeUsage(name string, types []string)
	GetHandlerForMap(action string, allHandlers map[string]Handler) (Handler, error)
	HideMovieRecommendation(client *internal.Client, options *str.Options) (*str.Response, error)
	HideShowRecommendation(client *internal.Client, options *str.Options) (*str.Response, error)
	Notes(client *internal.Client, notes *str.Notes, options *str.Options) (*str.Notes, *str.Response, error)
	PauseScrobble(client *internal.Client, scrobble *str.Scrobble, options *str.Options) (*str.Scrobble, *str.Response, error)
	ReadInput(options str.Options) (*str.ItemsList, error)
	Reply(client *internal.Client, id *int, comment *str.Comment, options *str.Options) (*str.Comment, *str.Response, error)
	StartScrobble(client *internal.Client, scrobble *str.Scrobble, options *str.Options) (*str.Scrobble, *str.Response, error)
	StopScrobble(client *internal.Client, scrobble *str.Scrobble, options *str.Options) (*str.Scrobble, *str.Response, error)
	ToTimestamp(at string) *str.Timestamp
	UpdateComment(client *internal.Client, options *str.Options, comment *str.Comment) (*str.Comment, *str.Response, error)
	UpdateHistoryListWithType(data []*str.ExportlistItem, strtype *string) []*str.ExportlistItem
	UpdateNotes(client *internal.Client, options *str.Options, notes *str.Notes) (*str.Notes, *str.Response, error)
	UsersAddToHiddenItems(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.AddResult, error)
	UsersRemoveHiddenItems(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.RemoveResult, error)
	ValidPrivacy(options *str.Options) error
	FetchUsersLikes(client *internal.Client, options *str.Options, page int) ([]*str.UserLike, error)
}

// CommonLogic struct for common methods
type CommonLogic struct{}

// Ensure CommonLogic implements CommonInterface at compile time
var _ CommonInterface = (*CommonLogic)(nil)

// CheckDates if dates are ok
func (*CommonLogic) CheckDates(from string, to string, tz string) error {
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

// CreateItemsToAddRatings helper function to create rating items from list
func (*CommonLogic) CreateItemsToAddRatings(items *str.ItemsList) str.RatingItems {
	movies := []str.Movie{}
	for _, m := range *items.Movies {
		movie := str.Movie{
			Title:     m.Title,
			Year:      m.Year,
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
			RatedAt:   m.RatedAt,
		}
		if m.Rating != nil {
			rating := float32(*m.Rating)
			movie.Rating = &rating
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
		if m.Rating != nil {
			rating := float32(*m.Rating)
			show.Rating = &rating
			show.RatedAt = m.RatedAt
		}
		if len(*m.Seasons) == consts.ZeroValue {
			show.Seasons = nil
		}
		shows = append(shows, show)
	}
	seasons := []str.Season{}
	for _, m := range *items.Seasons {
		season := str.Season{
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
			RatedAt:   m.RatedAt,
		}
		if m.Rating != nil {
			rating := float32(*m.Rating)
			season.Rating = &rating
		}

		seasons = append(seasons, season)
	}
	episodes := []str.Episode{}
	for _, m := range *items.Episodes {
		episode := str.Episode{
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
			RatedAt:   m.RatedAt,
		}
		if m.Rating != nil {
			rating := float32(*m.Rating)
			episode.Rating = &rating
		}

		episodes = append(episodes, episode)
	}

	return str.RatingItems{
		Movies:   &movies,
		Shows:    &shows,
		Seasons:  &seasons,
		Episodes: &episodes,
	}
}

// CreateItemsToRemoveRatings helper for items to remove from ratings
func (c *CommonLogic) CreateItemsToRemoveRatings(items *str.ItemsList) str.ItemsToRemove {
	return c.CreateItemsToRemove(items)
}

// OnlySeasonsIDs helper for list of season ids object
func (*CommonLogic) OnlySeasonsIDs(items *[]str.ExportlistItem) *[]str.Season {
	result := onlyIDs[str.Season](*items)
	return &result
}

// OnlyMoviesIDs helper for list of movies ids object
func (*CommonLogic) OnlyMoviesIDs(items *[]str.ExportlistItem) *[]str.Movie {
	result := onlyIDs[str.Movie](*items)
	return &result
}

// OnlyShowsIDs helper for list of shows ids object
func (*CommonLogic) OnlyShowsIDs(items *[]str.ExportlistItem) *[]str.Show {
	result := onlyIDs[str.Show](*items)
	return &result
}

// OnlyEpisodesIDs helper for list of episodes ids object
func (*CommonLogic) OnlyEpisodesIDs(items *[]str.ExportlistItem) *[]str.Episode {
	result := onlyIDs[str.Episode](*items)
	return &result
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

// CreateItemsToReorder helper to create list of watchlist ids to reorder
func (CommonLogic) CreateItemsToReorder(items *str.ItemsList) str.ItemsToReorder {
	reorder := []int64{}
	for _, m := range *items.Movies {
		reorder = append(reorder, *m.ID)
	}
	for _, m := range *items.Shows {
		reorder = append(reorder, *m.ID)
	}
	for _, m := range *items.Seasons {
		reorder = append(reorder, *m.ID)
	}
	for _, m := range *items.Episodes {
		reorder = append(reorder, *m.ID)
	}

	return str.ItemsToReorder{
		Rank: &reorder,
	}
}

// CreateItemsToAdd helper to create list of items to add to history
func (*CommonLogic) CreateItemsToAdd(items *str.ItemsList) str.HistoryItems {
	movies := []str.Movie{}
	for _, m := range *items.Movies {
		movie := str.Movie{
			Title:     m.Title,
			Year:      m.Year,
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
			HiddenAt:  m.HiddenAt,
			Notes:     m.Notes,
		}
		movies = append(movies, movie)
	}
	shows := []str.Show{}
	for _, m := range *items.Shows {
		show := str.Show{
			Title:    m.Title,
			Year:     m.Year,
			IDs:      m.IDs,
			Seasons:  m.Seasons,
			Notes:    m.Notes,
			HiddenAt: m.HiddenAt,
		}
		shows = append(shows, show)
	}
	seasons := []str.Season{}
	for _, m := range *items.Seasons {
		season := str.Season{
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
			HiddenAt:  m.HiddenAt,
			Notes:     m.Notes,
		}
		seasons = append(seasons, season)
	}
	episodes := []str.Episode{}
	for _, m := range *items.Episodes {
		episode := str.Episode{
			IDs:       m.IDs,
			WatchedAt: m.WatchedAt,
			HiddenAt:  m.HiddenAt,
			Notes:     m.Notes,
		}
		episodes = append(episodes, episode)
	}
	users := []str.UserProfile{}
	for _, m := range *items.Users {
		user := str.UserProfile{
			IDs:      m.IDs,
			HiddenAt: m.HiddenAt,
		}
		users = append(users, user)
	}

	return str.HistoryItems{
		Movies:   &movies,
		Shows:    &shows,
		Seasons:  &seasons,
		Episodes: &episodes,
		Users:    &users,
	}
}

// CreateCheckin helper function to create checkin object
func (c *CommonLogic) CreateCheckin(client *internal.Client, options *str.Options) (*str.Checkin, error) {
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
func (c *CommonLogic) CreateCheckinShowEpisode(client *internal.Client, options *str.Options) (*str.Checkin, error) {
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
func (c *CommonLogic) CreateScrobble(client *internal.Client, options *str.Options) (*str.Scrobble, error) {
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
func (c *CommonLogic) CreateScrobbleShowEpisode(client *internal.Client, options *str.Options) (*str.Scrobble, error) {
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
func (*CommonLogic) CheckTypes(options *str.Options) error {
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
func (*CommonLogic) ConvertDateString(dateStr string, outputFormat string, tz string, full bool) string {
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
func (*CommonLogic) ToTimestamp(at string) *str.Timestamp {
	// Parse the input date string using YYYY-MM-DD
	parsedDate, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return &str.Timestamp{}
	}

	return &str.Timestamp{Time: parsedDate.UTC()}
}

// CurrentDateString return current date string from user timezone
func (*CommonLogic) CurrentDateString(tz string, full bool) string {
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
func (*CommonLogic) DateLastDays(days int, tz string, full bool) string {
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

// ConvertBytesToItemsList convert bytes to struct
func (c *CommonLogic) ConvertBytesToItemsList(data []byte, action string, stype string) (*str.ItemsList, error) {
	var list []*str.ExportlistItem
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	items := c.InitItemsList()
	switch action {
	case consts.AddToHistory, consts.RemoveFromHistory, consts.AddToRatings, consts.RemoveFromRatings:
		items = c.ListToItemsAgregate(items, list, stype)
		return items.Uniq(), nil
	case consts.AddToCollection, consts.RemoveFromCollection, consts.RemoveFromWatchlist, consts.AddToWatchlist,
		consts.ReorderWatchlist, consts.AddToFavorites, consts.RemoveFromFavorites, consts.ReorderFavorites:
		items = c.ListToItemsCollection(items, list, stype)
		return items, nil
	case consts.AddHiddenItems, consts.RemoveHiddenItems:
		items = c.ListToItemsCollectionAgregate(items, list, stype)
		return items, nil
	default:
		return nil, errors.New(consts.UnknownItemsListType)
	}
}

// ListToItemsCollectionAgregate helper function to handle all types at once
func (c *CommonLogic) ListToItemsCollectionAgregate(items *str.ItemsList, list []*str.ExportlistItem, stype string) *str.ItemsList {
	if stype != "" {
		return c.ListToItemsCollection(items, list, stype)
	}

	out := c.InitItemsList()

	it := c.ListToItemsCollection(items, list, consts.Movie)
	for _, item := range *it.Movies {
		*out.Movies = append(*out.Movies, item)
	}
	it = c.ListToItemsCollection(items, list, consts.Show)
	for _, item := range *it.Shows {
		*out.Shows = append(*out.Shows, item)
	}
	it = c.ListToItemsCollection(items, list, consts.Season)
	for _, item := range *it.Seasons {
		*out.Seasons = append(*out.Seasons, item)
	}
	it = c.ListToItemsCollection(items, list, consts.Episode)
	for _, item := range *it.Episodes {
		*out.Episodes = append(*out.Episodes, item)
	}
	it = c.ListToItemsCollection(items, list, consts.User)
	for _, item := range *it.Users {
		*out.Users = append(*out.Users, item)
	}

	return out
}

// ListToItemsAgregate helper function to update itemslists depends on type: movies,shows,seasons,episodes,all
func (c *CommonLogic) ListToItemsAgregate(items *str.ItemsList, list []*str.ExportlistItem, stype string) *str.ItemsList {
	if stype != "all" {
		return c.ListToItems(items, list, stype)
	}

	out := c.InitItemsList()

	it := c.ListToItems(items, list, consts.Movies)
	for _, item := range *it.Movies {
		*out.Movies = append(*out.Movies, item)
	}
	it = c.ListToItems(items, list, consts.Shows)
	for _, item := range *it.Shows {
		*out.Shows = append(*out.Shows, item)
	}
	it = c.ListToItems(items, list, consts.Seasons)
	for _, item := range *it.Seasons {
		*out.Seasons = append(*out.Seasons, item)
	}
	it = c.ListToItems(items, list, consts.Episodes)
	for _, item := range *it.Episodes {
		*out.Episodes = append(*out.Episodes, item)
	}

	return out
}

// InitItemsList helper function to create InitItemsList with empty elements
func (*CommonLogic) InitItemsList() *str.ItemsList {
	list := new(str.ItemsList)
	list.Movies = &[]str.ExportlistItem{}
	list.Shows = &[]str.ExportlistItem{}
	list.Seasons = &[]str.ExportlistItem{}
	list.Episodes = &[]str.ExportlistItem{}
	list.Users = &[]str.ExportlistItem{}
	list.IDs = &[]int64{}
	return list
}

// ListToItemsCollection helper function to convert one list to another
func (*CommonLogic) ListToItemsCollection(items *str.ItemsList, list []*str.ExportlistItem, stype string) *str.ItemsList {
	for _, val := range list {
		if val.ID != nil {
			*items.IDs = append(*items.IDs, *val.ID)
		}

		if val.Movie != nil && isMovieType(stype) {
			e := str.ExportlistItem{}
			e.Title = val.Movie.Title
			e.Year = val.Movie.Year
			e.IDs = val.Movie.IDs
			e.UpdateCollectedData(val)
			e.ID = val.ID
			*items.Movies = append(*items.Movies, e)
		}
		if val.Show != nil && isShowType(stype) {
			e := str.ExportlistItem{}
			e.Title = val.Show.Title
			e.Year = val.Show.Year
			e.IDs = val.Show.IDs
			e.UpdateCollectedData(val)
			e.ID = val.ID
			*items.Shows = append(*items.Shows, e)
		}
		if val.Season != nil && isSeasonType(stype) {
			e := str.ExportlistItem{}
			e.IDs = val.Season.IDs
			val.Season.UpdateCollectedData(val)
			e.UpdateCollectedData(val)
			e.ID = val.ID
			*items.Seasons = append(*items.Seasons, e)
		}
		if val.Episode != nil && isEpisodeType(stype) {
			e := str.ExportlistItem{}
			e.IDs = val.Episode.IDs
			e.UpdateCollectedData(val)
			e.ID = val.ID
			*items.Episodes = append(*items.Episodes, e)
		}
	}

	return items
}

// ListToItems convert list to history items to struct
func (*CommonLogic) ListToItems(items *str.ItemsList, list []*str.ExportlistItem, stype string) *str.ItemsList {
	// Group by movie
	moviesMap := map[int64]str.OutputMovie{}
	// Group by show
	showsMap := map[int64]str.OutputShow{}
	// Group by season
	seasonsMap := map[int64]str.OutputSeason{}
	// Group by episodes
	episodesMap := map[int64]str.OutputEpisode{}

	for _, item := range list {
		switch stype {
		case consts.Movies:
			if item.Movie != nil {
				if item.ID != nil {
					*items.IDs = append(*items.IDs, *item.ID)
				}
				movieID := item.Movie.IDs.Trakt
				movie, exists := moviesMap[*movieID]
				if !exists {
					movie = str.OutputMovie{
						Title:     item.Movie.Title,
						Year:      item.Movie.Year,
						IDs:       item.Movie.IDs,
						WatchedAt: item.WatchedAt,
						RatedAt:   item.RatedAt,
						Rating:    item.Rating,
					}
				}

				moviesMap[*movieID] = movie
			}

		case consts.Shows:

			if item.Show != nil {
				showID := item.Show.IDs.Trakt
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
						RatedAt:   item.RatedAt,
						Rating:    item.Rating,
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

				if season == nil && item.Episode != nil {
					newSeason := str.Season{
						Number: item.Episode.Season,
					}
					*show.Seasons = append(*show.Seasons, newSeason)
					season = &(*show.Seasons)[len(*show.Seasons)-1]
					season.Episodes = &[]str.Episode{}
				}

				if item.Episode != nil {
					// Add episode
					newEpisode := &str.Episode{Number: item.Episode.Number, WatchedAt: item.WatchedAt}
					if item.Rating != nil {
						rating := float32(*item.Rating)
						newEpisode.Rating = &rating
						newEpisode.RatedAt = item.RatedAt
					}

					*season.Episodes = append(*season.Episodes, *newEpisode)
				}
				showsMap[showIDnp] = show
			}

		case consts.Seasons:
			if item.Season != nil {
				seasonID := item.Season.IDs.Trakt
				season, exists := seasonsMap[*seasonID]
				if !exists {
					season = str.OutputSeason{
						Title:     item.Season.Title,
						IDs:       item.Season.IDs,
						WatchedAt: item.WatchedAt,
						RatedAt:   item.RatedAt,
						Rating:    item.Rating,
					}
				}
				seasonsMap[*seasonID] = season
			}
		case consts.Episodes:
			if item.Episode != nil {
				episodeID := item.Episode.IDs.Trakt
				episode, exists := episodesMap[*episodeID]
				if !exists {
					episode = str.OutputEpisode{
						Title:     item.Episode.Title,
						IDs:       item.Episode.IDs,
						WatchedAt: item.WatchedAt,
						RatedAt:   item.RatedAt,
						Rating:    item.Rating,
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
			RatedAt:   s.RatedAt,
			Rating:    s.Rating,
			IDs:       s.IDs,
		})
	}
	for _, s := range seasonsMap {
		*items.Seasons = append(*items.Seasons, str.ExportlistItem{
			Season:    &str.Season{Title: s.Title, IDs: s.IDs},
			WatchedAt: s.WatchedAt,
			RatedAt:   s.RatedAt,
			Rating:    s.Rating,
			IDs:       s.IDs,
		})
	}
	for _, s := range moviesMap {
		*items.Movies = append(*items.Movies, str.ExportlistItem{
			Movie:     &str.Movie{Title: s.Title, Year: s.Year, IDs: s.IDs},
			WatchedAt: s.WatchedAt,
			RatedAt:   s.RatedAt,
			Rating:    s.Rating,
			IDs:       s.IDs,
		})
	}

	for _, s := range showsMap {
		*items.Shows = append(*items.Shows, str.ExportlistItem{
			Show:      &str.Show{Title: s.Title, Year: s.Year, IDs: s.IDs},
			Seasons:   s.Seasons,
			WatchedAt: s.WatchedAt,
			RatedAt:   s.RatedAt,
			Rating:    s.Rating,
			IDs:       s.IDs,
		})
	}

	return items
}

// ReadInput read data from stdin or from file
func (c *CommonLogic) ReadInput(options str.Options) (*str.ItemsList, error) {
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

// FetchRatings returns users ratings filtered by type.
func (c CommonLogic) FetchRatings(client *internal.Client, options *str.Options, page int) ([]*str.RatingListItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	r := options.Rating.String()
	list, resp, err := client.Sync.GetRatings(
		client.BuildCtxFromOptions(options),
		&options.Type,
		&r,
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
		nextPageItems, err := c.FetchRatings(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// UpdateHistoryListWithType helper function to update History list with new type
func (*CommonLogic) UpdateHistoryListWithType(data []*str.ExportlistItem, strtype *string) []*str.ExportlistItem {
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

// FetchWatchlist helper function to fetch user watchlist
func (c *CommonLogic) FetchWatchlist(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetWatchlist(
		client.BuildCtxFromOptions(options),
		&options.Type,
		&options.SortBy,
		&options.SortHow,
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
		nextPageItems, err := c.FetchWatchlist(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// FetchFavorites helper function to fetch user favorites
func (c CommonLogic) FetchFavorites(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetFavorites(
		client.BuildCtxFromOptions(options),
		&options.Type,
		&options.SortBy,
		&options.SortHow,
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
		nextPageItems, err := c.FetchFavorites(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// FetchPendingFollowingRequests helper function to fetch pending following requests
func (CommonLogic) FetchPendingFollowingRequests(client *internal.Client, options *str.Options) ([]*str.FollowRequest, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Users.GetPendingFollowingRequests(
		client.BuildCtxFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// FetchFollowRequests helper function to fetch follow requests
func (CommonLogic) FetchFollowRequests(client *internal.Client, options *str.Options) ([]*str.FollowRequest, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	list, _, err := client.Users.GetFollowRequests(
		client.BuildCtxFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}

// ApproveFollowRequest helper function to approve follow request
func (CommonLogic) ApproveFollowRequest(client *internal.Client, options *str.Options) (*str.FollowRequest, *str.Response, error) {
	result, resp, err := client.Users.ApproveFollowRequest(
		client.BuildCtxFromOptions(options),
		options.FollowerRequest,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

// DenyFollowRequest helper function to deny follow request
func (CommonLogic) DenyFollowRequest(client *internal.Client, options *str.Options) (*str.FollowRequest, *str.Response, error) {
	result, resp, err := client.Users.DenyFollowRequest(
		client.BuildCtxFromOptions(options),
		options.FollowerRequest,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

// FetchUsersHiddenItems helper function to fetch users:hidden items
func (c *CommonLogic) FetchUsersHiddenItems(client *internal.Client, options *str.Options, page int) ([]*str.HiddenItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, Type: options.Type}
	list, resp, err := client.Users.GetHiddenItems(
		client.BuildCtxFromOptions(options),
		&options.Section,
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
		nextPageItems, err := c.FetchUsersHiddenItems(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// CreateItemsToHidden helper function to prepare items to hidden items
func (c CommonLogic) CreateItemsToHidden(section string, items *str.ItemsList) str.HistoryItems {
	i := c.CreateItemsToAdd(items)
	switch section {
	case consts.Calendar:
		return str.HistoryItems{Movies: i.Movies, Shows: i.Shows}
	case consts.ProgressWatched, consts.ProgressCollected:
		return str.HistoryItems{Shows: i.Shows, Seasons: i.Seasons}
	case consts.Recommendations:
		return str.HistoryItems{Movies: i.Movies, Shows: i.Shows}
	case consts.Comments:
		return str.HistoryItems{Users: i.Users}
	case consts.Dropped:
		return str.HistoryItems{Shows: i.Shows}
	default:
		return str.HistoryItems{}
	}
}

// UsersAddToHiddenItems helper function to users: add hidden items
func (CommonLogic) UsersAddToHiddenItems(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.AddResult, error) {
	result, err := client.Users.AddHiddenItems(
		client.BuildCtxFromOptions(options),
		items,
		options.Section,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UsersRemoveHiddenItems helper function to users: remove hidden items
func (CommonLogic) UsersRemoveHiddenItems(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.RemoveResult, error) {
	result, err := client.Users.RemoveHiddenItems(
		client.BuildCtxFromOptions(options),
		items,
		options.Section,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FetchUsersLikes helper function to users: likes
func (c CommonLogic) FetchUsersLikes(client *internal.Client, options *str.Options, page int) ([]*str.UserLike, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Users.GetLikes(
		client.BuildCtxFromOptions(options),
		&options.UserName,
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
		nextPageItems, err := c.FetchUsersLikes(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
