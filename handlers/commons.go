// Package handlers used to handle module actions
package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// CommonInterface interface
type CommonInterface interface {
	FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, error)
	FetchShow(client *internal.Client, id *int) (*str.Show, error)
	FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error)
	FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error)
	CheckSeasonNumber(code *string) (*string, *string, error)
	Checkin(client *internal.Client, checkin *str.CheckIn) (*str.CheckIn, *str.Response, error)
	Comment(client *internal.Client, checkin *str.Comment) (*str.Comment, *str.Response, error)
}

// CommonLogic struct for common methods
type CommonLogic struct{}

// FetchMovie helper function to fetch movie object
func (*CommonLogic) FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, error) {
	movieID := options.TraktID
	result, _, err := client.Movies.GetMovie(
		context.Background(),
		&movieID,
	)

	return result, err
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

// FetchEpisode helper function to fetch episode object
func (*CommonLogic) FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error) {
	episodeID := options.TraktID
	result, _, err := client.Episodes.GetEpisode(
		context.Background(),
		&episodeID,
	)

	return result, err
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
