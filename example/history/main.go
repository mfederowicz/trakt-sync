// Package main prints the authenticated user's watch history, page by page, in their local timezone:
//
//	TRAKT_CLIENT_ID=... TRAKT_ACCESS_TOKEN=... go run ./history
//
// Get an access token with the devicecode example.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// maxPages keeps the example short; drop it to fetch the whole history.
const maxPages = 3

func main() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithAuthToken(os.Getenv("TRAKT_ACCESS_TOKEN")).
		WithUserAgent("trakt-sync-example/1.0")

	// timestamps in responses are converted to this location
	ctx := trakt.WithTimezone(context.Background(), time.Local)

	opts := &uri.ListOptions{Page: 1, Limit: 20}
	for {
		items, resp, err := client.Sync.GetWatchedHistory(ctx, 0, "", opts) // 0 and "": all entries of all types
		if err != nil {
			log.Fatalf("history page %d: %v", opts.Page, err)
		}
		for _, item := range items {
			fmt.Println(describe(item))
		}
		if !client.HavePages(opts.Page, resp, maxPages) {
			break
		}
		opts.Page++
	}
}

func describe(item *str.ExportlistItem) string {
	watched := "?"
	if item.WatchedAt != nil {
		watched = item.WatchedAt.Format("2006-01-02 15:04")
	}
	switch {
	case item.Movie != nil && item.Movie.Title != nil:
		return fmt.Sprintf("%s  movie    %s", watched, *item.Movie.Title)
	case item.Show != nil && item.Show.Title != nil && item.Episode != nil && item.Episode.Season != nil && item.Episode.Number != nil:
		return fmt.Sprintf("%s  episode  %s %dx%02d", watched, *item.Show.Title, *item.Episode.Season, *item.Episode.Number)
	default:
		return fmt.Sprintf("%s  %s", watched, "(unknown item)")
	}
}
