// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MediaPopularHandler struct for handler
type MediaPopularHandler struct{}

// Handle to handle media: popular action
func (MediaPopularHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns popular movies and shows.")
	return exportMedia(client, options, func(opts *uri.ListOptions) ([]*str.Media, *str.Response, error) {
		return client.Media.GetPopularMedia(client.BuildCtxFromOptions(options), opts)
	})
}
