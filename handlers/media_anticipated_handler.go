// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MediaAnticipatedHandler struct for handler
type MediaAnticipatedHandler struct{}

// Handle to handle media: anticipated action
func (MediaAnticipatedHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns anticipated movies and shows.")
	return exportMedia(client, options, func(opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error) {
		return client.Media.GetAnticipatedMedia(client.BuildCtxFromOptions(options), opts)
	})
}
