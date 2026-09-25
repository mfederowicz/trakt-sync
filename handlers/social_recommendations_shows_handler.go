// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SocialRecommendationsShowsHandler struct for handler
type SocialRecommendationsShowsHandler struct{}

// Handle to handle social_recommendations: shows action
func (SocialRecommendationsShowsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns show recommendations based on the people you follow.")
	return exportSocialRecommendations(client, options, func(opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
		return client.SocialRecommendations.GetSocialShowRecommendations(client.BuildCtxFromOptions(options), opts)
	})
}
