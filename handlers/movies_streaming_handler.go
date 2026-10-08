// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MoviesStreamingHandler struct for handler
type MoviesStreamingHandler struct{}

// Handle to handle movies: streaming action
func (MoviesStreamingHandler) Handle(options *str.Options, client *trakt.Client) error {
	periods := cfg.ModuleActionConfig[consts.Movies+":"+consts.Streaming].Period
	if !cfg.IsValidConfigType(periods, options.Period) {
		return fmt.Errorf("period '%s' is not valid for streaming, avaliable periods:%s", options.Period, periods)
	}

	if err := checkMovieFilters(options); err != nil {
		return err
	}

	printer.Println("Returns the most streamed movies in the specified time period.")
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
		opts.Filters = mediaFilters(options)
		opts.Status = movieStatus(options)
		return client.Movies.GetStreamingMovies(cli.ContextFromOptions(options), options.Period, opts)
	})
	if err != nil {
		return endpointNotLiveError(err)
	}

	return writeMoviesItems(options, result)
}
