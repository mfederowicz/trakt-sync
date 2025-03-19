// Package handlers used to handle module actions
package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// CommonInterface interface
type CommonInterface interface {
	FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, error)
	FetchShow(client *internal.Client, id *int) (*str.Show, error)
	FetchSeason(client *internal.Client, id *int) (*str.Season, error)
	FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error)
	FetchList(client *internal.Client, options *str.Options) (*str.PersonalList, error)
	FetchComment(client *internal.Client, options *str.Options) (*str.Comment, error)
	FetchCommentItem(client *internal.Client, options *str.Options) (*str.CommentMediaItem, error)
	FetchCommentUserLikes(client *internal.Client, options *str.Options) (*str.CommentUserLike, error)
	FetchTrendingComments(client *internal.Client, options *str.Options) (*str.CommentItem, error)
	FetchRecentComments(client *internal.Client, options *str.Options) (*str.CommentItem, error)
	FetchUpdatedComments(client *internal.Client, options *str.Options) (*str.CommentItem, error)
	UpdateComment(client *internal.Client, options *str.Options) (*str.Comment, error)
	DeleteComment(client *internal.Client, options *str.Options) (*str.Comment, *str.Response, error)
	FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error)
	CheckSeasonNumber(code *string) (*string, *string, error)
	Checkin(client *internal.Client, checkin *str.CheckIn) (*str.CheckIn, *str.Response, error)
	Comment(client *internal.Client, comment *str.Comment) (*str.Comment, *str.Response, error)
	Reply(client *internal.Client, id *int, comment *str.Comment) (*str.Comment, *str.Response, error)
}

// CommonLogic struct for common methods
type CommonLogic struct{}

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
	showID := options.TraktID

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
	seasonID := options.TraktID
	result, _, err := client.Seasons.GetSeason(
		context.Background(),
		&seasonID,
		&opts,
	)

	return result, err
}

// FetchEpisode helper function to fetch episode object
func (*CommonLogic) FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error) {
	episodeID := options.TraktID
	result, _, err := client.Episodes.GetEpisode(
		context.Background(),
		&episodeID,
	)

	return result, err
}

// FetchList helper function to fetch list object
func (*CommonLogic) FetchList(client *internal.Client, options *str.Options) (*str.PersonalList, error) {
	listID := options.TraktID
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

// DeleteComment helper function to delete comment object
func (*CommonLogic) DeleteComment(client *internal.Client, options *str.Options) (*str.Response, error) {
	commentID := options.CommentID
	resp, err := client.Comments.DeleteComment(
		context.Background(),
		&commentID,
	)

	return resp, err
}

// FetchUserConnections helper function to fetch connections object
func (*CommonLogic) FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error) {
	result, _, err := client.Users.RetrieveSettings(
		context.Background(),
	)

	return result.Connections, err
}

// Checkin helper function to post checkin object
func (*CommonLogic) Checkin(client *internal.Client, checkin *str.CheckIn) (*str.CheckIn, *str.Response, error) {
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
		return nil, nil, errors.New("invalid episode_code format")
	}

	if parts := strings.Split(code, "x"); len(parts) == consts.TwoValue {
		season, _ := strconv.Atoi(parts[0])
		episode, _ := strconv.Atoi(parts[1])
		return &season, &episode, nil
	}

	return nil, nil, errors.New("invalid episode_code format")
}
