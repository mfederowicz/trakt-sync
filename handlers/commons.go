// Package handlers used to handle module actions
package handlers

import (
	"context"
	"errors"
	"fmt"
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
	ToTimestamp(at string) (*str.Timestamp, error)
	ConvertDateString(date string, out string) string
}

// CommonLogic struct for common methods
type CommonLogic struct{}

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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
		&seasonID,
		&opts,
	)

	return result, err
}

// FetchEpisode helper function to fetch episode object
func (*CommonLogic) FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error) {
	episodeID := options.InternalID
	result, _, err := client.Episodes.GetEpisode(
		context.Background(),
		&episodeID,
	)

	return result, err
}

// FetchPerson helper function to fetch person object
func (*CommonLogic) FetchPerson(client *internal.Client, options *str.Options) (*str.Person, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	personID := options.InternalID
	result, _, err := client.People.GetSinglePerson(
		context.Background(),
		&personID,
		&opts,
	)

	return result, err
}

// FetchList helper function to fetch list object
func (*CommonLogic) FetchList(client *internal.Client, options *str.Options) (*str.PersonalList, error) {
	listID := options.InternalID
	result, _, err := client.Lists.GetList(
		context.Background(),
		&listID,
	)

	return result, err
}

// FetchComment helper function to fetch comment object
func (*CommonLogic) FetchComment(client *internal.Client, options *str.Options) (*str.Comment, error) {
	commentID := options.CommentID
	result, _, err := client.Comments.GetComment(
		context.Background(),
		&commentID,
	)

	return result, err
}

// FetchNotes helper function to fetch notes object
func (*CommonLogic) FetchNotes(client *internal.Client, options *str.Options) (*str.Notes, error) {
	notesID := options.InternalID
	result, _, err := client.Notes.GetNotes(
		context.Background(),
		&notesID,
	)

	return result, err
}

// FetchNotesItem helper function to fetch notes attached item object
func (*CommonLogic) FetchNotesItem(client *internal.Client, options *str.Options) (*str.NotesItem, error) {
	notesID := options.InternalID
	result, _, err := client.Notes.GetNotesItem(
		context.Background(),
		&notesID,
	)

	return result, err
}

// FetchCommentItem helper function to fetch comment media item object
func (*CommonLogic) FetchCommentItem(client *internal.Client, options *str.Options) (*str.CommentMediaItem, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	commentID := options.CommentID
	result, _, err := client.Comments.GetCommentItem(
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
		&commentID,
		comment,
	)

	return result, resp, err
}

// UpdateNotes helper function to put notes object
func (*CommonLogic) UpdateNotes(client *internal.Client, options *str.Options, notes *str.Notes) (*str.Notes, *str.Response, error) {
	notesID := options.InternalID
	result, resp, err := client.Notes.UpdateNotes(
		context.Background(),
		&notesID,
		notes,
	)

	return result, resp, err
}

// DeleteComment helper function to delete comment object
func (*CommonLogic) DeleteComment(client *internal.Client, options *str.Options) (*str.Response, error) {
	commentID := options.CommentID
	resp, err := client.Comments.DeleteComment(
		context.Background(),
		&commentID,
	)

	return resp, err
}

// DeleteNotes helper function to delete notes object
func (*CommonLogic) DeleteNotes(client *internal.Client, options *str.Options) (*str.Response, error) {
	notesID := options.InternalID
	resp, err := client.Notes.DeleteNotes(
		context.Background(),
		&notesID,
	)

	return resp, err
}

// HideMovieRecommendation helper function to hide movie recommendations
func (*CommonLogic) HideMovieRecommendation(client *internal.Client, options *str.Options) (*str.Response, error) {
	movieID := options.InternalID
	resp, err := client.Recommendations.HideMovieRecommendation(
		context.Background(),
		&movieID,
	)

	return resp, err
}

// HideShowRecommendation helper function to hide show recommendations
func (*CommonLogic) HideShowRecommendation(client *internal.Client, options *str.Options) (*str.Response, error) {
	showID := options.InternalID
	resp, err := client.Recommendations.HideShowRecommendation(
		context.Background(),
		&showID,
	)

	return resp, err
}

// FetchUserConnections helper function to fetch connections object
func (*CommonLogic) FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error) {
	result, _, err := client.Users.RetrieveSettings(
		context.Background(),
	)
	if err != nil {
		return nil, fmt.Errorf(consts.UserSettingsError, err)
	}

	return result.Connections, err
}

// StartScrobble helper function to start scrobble
func (*CommonLogic) StartScrobble(client *internal.Client, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error) {
	result, resp, err := client.Scrobble.StartScrobble(
		context.Background(),
		scrobble,
	)

	return result, resp, err
}

// StopScrobble helper function to stop scrobble
func (*CommonLogic) StopScrobble(client *internal.Client, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error) {
	result, resp, err := client.Scrobble.StopScrobble(
		context.Background(),
		scrobble,
	)

	return result, resp, err
}

// PauseScrobble helper function to pause scrobble
func (*CommonLogic) PauseScrobble(client *internal.Client, scrobble *str.Scrobble) (*str.Scrobble, *str.Response, error) {
	result, resp, err := client.Scrobble.PauseScrobble(
		context.Background(),
		scrobble,
	)

	return result, resp, err
}

// Checkin helper function to post checkin object
func (*CommonLogic) Checkin(client *internal.Client, checkin *str.Checkin) (*str.Checkin, *str.Response, error) {
	result, resp, err := client.Checkin.CheckintoAnItem(
		context.Background(),
		checkin,
	)

	return result, resp, err
}

// Comment helper function to post comment object
func (*CommonLogic) Comment(client *internal.Client, comment *str.Comment) (*str.Comment, *str.Response, error) {
	result, resp, err := client.Comments.PostAComment(
		context.Background(),
		comment,
	)
	return result, resp, err
}

// Notes helper function to post notes object
func (*CommonLogic) Notes(client *internal.Client, notes *str.Notes) (*str.Notes, *str.Response, error) {
	result, resp, err := client.Notes.AddNotes(
		context.Background(),
		notes,
	)
	return result, resp, err
}

// Reply helper function to post reply object
func (*CommonLogic) Reply(client *internal.Client, id *int, reply *str.Comment) (*str.Comment, *str.Response, error) {
	result, resp, err := client.Comments.ReplyAComment(
		context.Background(),
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
func (CommonLogic) ConvertDateString(dateStr string, outputFormat string, tz string) string {
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
