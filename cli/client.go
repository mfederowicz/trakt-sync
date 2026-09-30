// Package cli for basic cli functions
package cli

import (
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// NewClient builds the API client for the CLI: client id from the config, the trakt-sync User-Agent,
// and the access token when there is one (an empty token is not sent).
func NewClient(config *cfg.Config, token str.Token) *trakt.Client {
	return trakt.NewClient(nil).
		WithClientID(config.ClientID).
		WithUserAgent(UserAgent()).
		WithAuthToken(token.AccessToken)
}
