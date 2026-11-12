// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MoviesAnticipatedRecommendationsHandler struct for handler
type MoviesAnticipatedRecommendationsHandler struct{}

// Handle to handle movies: anticipated_recommendations action
func (MoviesAnticipatedRecommendationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	return exportAnticipatedRecommendations(client, options, func(opts *uri.ListOptions) ([]*str.AnticipatedRecommendation, *str.Response, error) {
		return client.Movies.GetAnticipatedRecommendations(cli.ContextFromOptions(options), opts)
	})
}
