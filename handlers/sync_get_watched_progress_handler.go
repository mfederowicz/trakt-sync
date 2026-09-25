// Package handlers used to handle module actions
package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SyncGetWatchedProgressHandler struct for handler
type SyncGetWatchedProgressHandler struct{}

// Handle to handle sync: get_watched_progress action
func (SyncGetWatchedProgressHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.HideCompleted && options.HideNotCompleted {
		return errors.New(consts.HideBothProgressMsg)
	}

	printer.Println("Get watched progress")
	opts := uri.SyncProgressOptions{
		Extended:         options.ExtendedInfo,
		Limit:            options.PerPage,
		SortBy:           options.SortBy,
		SortHow:          options.SortHow,
		LifetimeStats:    options.LifetimeStats,
		HideCompleted:    options.HideCompleted,
		HideNotCompleted: options.HideNotCompleted,
		OnlyRewatching:   options.OnlyRewatching,
	}
	result, err := fetchSyncProgress(client, options, consts.DefaultPage, func(ctx context.Context, page int) ([]*str.ShowProgress, *str.Response, error) {
		opts.Page = page
		return client.Sync.GetWatchedProgress(ctx, &opts)
	})
	if err != nil {
		return fmt.Errorf("get watched progress error: %w", err)
	}

	return writeSyncProgress(options, result)
}
