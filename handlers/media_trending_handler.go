// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MediaTrendingHandler struct for handler
type MediaTrendingHandler struct{}

// Handle to handle media: trending action
func (MediaTrendingHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns trending movies and shows.")
	return exportMedia(client, options, func(opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error) {
		return client.Media.GetTrendingMedia(client.BuildCtxFromOptions(options), opts)
	})
}
