// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MoviesHotHandler struct for handler
type MoviesHotHandler struct{}

// Handle to handle movies: hot action
func (MoviesHotHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns hot movies, based on current list activity.")
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
		return client.Movies.GetHotMovies(client.BuildCtxFromOptions(options), opts)
	})
	if err != nil {
		return endpointNotLiveError(err)
	}

	return writeMoviesItems(options, result)
}
