// Package main lists the 10 trending movies. It needs only an API app client id:
//
//	TRAKT_CLIENT_ID=... go run ./trending
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

func main() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("trakt-sync-example/1.0")

	movies, resp, err := client.Movies.GetTrendingMovies(context.Background(), &uri.ListOptions{Limit: 10})
	if err != nil {
		log.Fatalf("trending movies: %v", err)
	}

	for i, item := range movies {
		if item.Movie == nil || item.Movie.Title == nil {
			continue
		}
		year := 0
		if item.Movie.Year != nil {
			year = *item.Movie.Year
		}
		watchers := 0
		if item.WatcherCount != nil {
			watchers = *item.WatcherCount
		}
		fmt.Printf("%2d. %s (%d), %d watching\n", i+1, *item.Movie.Title, year, watchers)
	}
	fmt.Printf("rate limit: %d of %d requests left\n", resp.Rate.Remaining, resp.Rate.Limit)
}
