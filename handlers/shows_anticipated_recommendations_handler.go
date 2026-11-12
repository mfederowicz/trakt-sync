// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// ShowsAnticipatedRecommendationsHandler struct for handler
type ShowsAnticipatedRecommendationsHandler struct{}

// Handle to handle shows: anticipated_recommendations action
func (ShowsAnticipatedRecommendationsHandler) Handle(options *str.Options, client *trakt.Client) error {
	return exportAnticipatedRecommendations(client, options, func(opts *uri.ListOptions) ([]*str.AnticipatedRecommendation, *str.Response, error) {
		return client.Shows.GetAnticipatedRecommendations(cli.ContextFromOptions(options), opts)
	})
}
