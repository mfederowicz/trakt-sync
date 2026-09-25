// Package handlers used to handle module actions
package handlers

import (
	"context"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SyncGetUpNextNitroHandler struct for handler
type SyncGetUpNextNitroHandler struct{}

// Handle to handle sync: get_up_next_nitro action
func (SyncGetUpNextNitroHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.Intent) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.SyncUpNextIntents, options.Intent) {
		return fmt.Errorf("intent '%s' is not valid, avaliable values: %v", options.Intent, cfg.SyncUpNextIntents)
	}
	if len(options.WatchNow) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.WatchNowFilters, options.WatchNow) {
		return fmt.Errorf("watchnow '%s' is not valid, avaliable values: %v", options.WatchNow, cfg.WatchNowFilters)
	}

	printer.Println("Get up next nitro progress")
	opts := uri.UpNextNitroOptions{
		Limit:          options.PerPage,
		SortBy:         options.SortBy,
		SortHow:        options.SortHow,
		Intent:         options.Intent,
		WatchNow:       options.WatchNow,
		Genres:         options.Genres,
		Subgenres:      options.Subgenres,
		Years:          options.Years,
		Ratings:        options.Ratings,
		StartDate:      options.MediaStartDate,
		EndDate:        options.MediaEndDate,
		Runtimes:       options.Runtimes,
		Countries:      options.Countries,
		Certifications: options.Certifications,
	}
	result, err := fetchSyncProgress(client, options, consts.DefaultPage, func(ctx context.Context, page int) ([]*str.ShowProgress, *str.Response, error) {
		opts.Page = page
		return client.Sync.GetUpNextNitro(ctx, &opts)
	})
	if err != nil {
		return fmt.Errorf("get up next nitro error: %w", err)
	}

	return writeSyncProgress(options, result)
}
