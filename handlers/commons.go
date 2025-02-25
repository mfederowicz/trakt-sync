// Package handlers used to handle module actions
package handlers

import (
	"context"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommonInterface interface
type CommonInterface interface {
	FetchMovie(client *internal.Client, options *str.Options) (*str.Movie, error)
	FetchEpisode(client *internal.Client, options *str.Options) (*str.Episode, error)
	FetchUserConnections(client *internal.Client, _ *str.Options) (*str.Connections, error)
	Checkin(client *internal.Client, checkin *str.CheckIn) (*str.CheckIn, *str.Response, error)
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
