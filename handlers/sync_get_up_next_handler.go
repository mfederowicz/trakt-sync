// Package handlers used to handle module actions
package handlers

import (
	"context"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SyncGetUpNextHandler struct for handler
type SyncGetUpNextHandler struct{}

// Handle to handle sync: get_up_next action
func (SyncGetUpNextHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get up next progress")
	opts := uri.SyncProgressOptions{
		Extended:      options.ExtendedInfo,
		Limit:         options.PerPage,
		SortBy:        options.SortBy,
		SortHow:       options.SortHow,
		IncludeStats:  options.IncludeStats,
		LifetimeStats: options.LifetimeStats,
	}
	result, err := fetchSyncProgress(client, options, consts.DefaultPage, func(ctx context.Context, page int) ([]*str.ShowProgress, *str.Response, error) {
		opts.Page = page
		return client.Sync.GetUpNext(ctx, &opts)
	})
	if err != nil {
		return fmt.Errorf("get up next error: %w", err)
	}

	return writeSyncProgress(options, result)
}
