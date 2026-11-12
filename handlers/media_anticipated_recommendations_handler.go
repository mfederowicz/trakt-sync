// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MediaAnticipatedRecommendationsHandler struct for handler
type MediaAnticipatedRecommendationsHandler struct{}

// Handle to handle media: anticipated_recommendations action
func (MediaAnticipatedRecommendationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	return exportAnticipatedRecommendations(client, options, func(opts *uri.ListOptions) ([]*str.AnticipatedRecommendation, *str.Response, error) {
		return client.Media.GetAnticipatedRecommendations(cli.ContextFromOptions(options), opts)
	})
}
