// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SocialRecommendationsMoviesHandler struct for handler
type SocialRecommendationsMoviesHandler struct{}

// Handle to handle social_recommendations: movies action
func (SocialRecommendationsMoviesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns movie recommendations based on the people you follow.")
	return exportSocialRecommendations(client, options, func(opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
		return client.SocialRecommendations.GetSocialMovieRecommendations(client.BuildCtxFromOptions(options), opts)
	})
}
