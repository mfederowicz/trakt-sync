// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MoviesStreamingHandler struct for handler
type MoviesStreamingHandler struct{}

// Handle to handle movies: streaming action
func (MoviesStreamingHandler) Handle(options *str.Options, client *internal.Client) error {
	periods := cfg.ModuleActionConfig[consts.Movies+":"+consts.Streaming].Period
	if !cfg.IsValidConfigType(periods, options.Period) {
		return fmt.Errorf("period '%s' is not valid for streaming, avaliable periods:%s", options.Period, periods)
	}

	printer.Println("Returns the most streamed movies in the specified time period.")
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
		return client.Movies.GetStreamingMovies(client.BuildCtxFromOptions(options), &options.Period, opts)
	})
	if err != nil {
		return err
	}

	return writeMoviesItems(options, result)
}
