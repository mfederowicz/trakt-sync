// Package trakt_test holds the documentation examples of the trakt package.
package trakt_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// A client needs the client id of your Trakt API app and a User-Agent that names your app.
func Example() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("my-app/1.0")

	movies, _, err := client.Movies.GetTrendingMovies(context.Background(), &uri.ListOptions{Limit: 10})
	if err != nil {
		fmt.Println("trending:", err)
		return
	}
	for _, item := range movies {
		if item.Movie != nil && item.Movie.Title != nil {
			fmt.Println(*item.Movie.Title)
		}
	}
}

// Routes for the authenticated user need an OAuth access token (see example/devicecode in the repository).
func ExampleClient_WithAuthToken() {
	base := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("my-app/1.0")

	// The With methods return copies, so one base client can serve several users.
	client := base.WithAuthToken(os.Getenv("TRAKT_ACCESS_TOKEN"))

	settings, _, err := client.Users.GetSettings(context.Background())
	if err != nil {
		fmt.Println("settings:", err)
		return
	}
	if settings.User != nil && settings.User.Username != nil {
		fmt.Println("logged in as", *settings.User.Username)
	}
}

// HavePages reads the pagination headers of a response; the last argument caps the pages (0 means no cap).
func ExampleClient_HavePages() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("my-app/1.0").
		WithAuthToken(os.Getenv("TRAKT_ACCESS_TOKEN"))
	ctx := context.Background()

	total := 0
	opts := &uri.ListOptions{Page: 1, Limit: 100}
	for {
		// 0 and "": all history entries of all types
		items, resp, err := client.Sync.GetWatchedHistory(ctx, 0, "", opts)
		if err != nil {
			fmt.Println("history:", err)
			return
		}
		total += len(items)
		if !client.HavePages(opts.Page, resp, 0) {
			break
		}
		opts.Page++
	}
	fmt.Println("history entries:", total)
}

// Common API errors have their own types; errors.As picks them out.
func ExampleNotFoundError() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("my-app/1.0")

	_, _, err := client.Movies.GetMovie(context.Background(), "no-such-movie-slug", &uri.ListOptions{})

	var notFound *trakt.NotFoundError
	var rateLimited *trakt.AbuseRateLimitError
	switch {
	case errors.As(err, &notFound):
		fmt.Println("no such movie")
	case errors.As(err, &rateLimited) && rateLimited.RetryAfter != nil:
		fmt.Println("rate limited, retry after", *rateLimited.RetryAfter)
	case err != nil:
		fmt.Println("request failed:", err)
	}
}

// Timestamps in responses are UTC unless the context carries another timezone.
func ExampleWithTimezone() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("my-app/1.0")

	ctx := trakt.WithTimezone(context.Background(), time.Local)
	movies, _, err := client.Movies.GetRecentlyUpdatedMovies(ctx, time.Now().AddDate(0, 0, -1).Format(time.DateOnly), &uri.ListOptions{Limit: 5})
	if err != nil {
		fmt.Println("updates:", err)
		return
	}
	for _, item := range movies {
		if item.UpdatedAt != nil {
			fmt.Println(item.UpdatedAt.Time) // in time.Local
		}
	}
}
