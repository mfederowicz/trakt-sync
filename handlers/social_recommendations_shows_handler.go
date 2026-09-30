// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SocialRecommendationsShowsHandler struct for handler
type SocialRecommendationsShowsHandler struct{}

// Handle to handle social_recommendations: shows action
func (SocialRecommendationsShowsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns show recommendations based on the people you follow.")
	return exportSocialRecommendations(client, options, func(opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
		return client.SocialRecommendations.GetSocialShowRecommendations(cli.ContextFromOptions(options), opts)
	})
}
