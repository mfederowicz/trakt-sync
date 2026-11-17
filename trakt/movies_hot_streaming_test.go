// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// TestMoviesServiceHotAndStreaming checks that the hot and streaming routes are locked in the registry: the methods
// return an *EndpointRemovedError and send no request.
func TestMoviesServiceHotAndStreaming(t *testing.T) {
	tests := []struct {
		name string
		want string
		call func(s *MoviesService, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error)
	}{
		{
			name: "hot",
			want: "GET movies/hot is locked since 2026-11-17: Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent",
			call: func(s *MoviesService, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
				return s.GetHotMovies(context.Background(), opts)
			},
		},
		{
			name: "streaming",
			want: "GET movies/streaming/{period} is locked since 2026-11-17: Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent",
			call: func(s *MoviesService, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
				return s.GetStreamingMovies(context.Background(), "daily", opts)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			var requests atomic.Int32
			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				requests.Add(1)
				w.WriteHeader(http.StatusNotFound)
			})

			got, resp, err := tt.call(setup.Client.Movies, &uri.ListOptions{Page: 2})
			if got != nil || resp != nil {
				t.Errorf("result is %v, response %v, want nil", got, resp)
			}
			removed := new(EndpointRemovedError)
			if !errors.As(err, &removed) {
				t.Fatalf("error is %v, want *EndpointRemovedError", err)
			}
			if err.Error() != tt.want {
				t.Errorf("error is %q, want %q", err.Error(), tt.want)
			}
			if n := requests.Load(); n != 0 {
				t.Errorf("%d requests sent, want 0", n)
			}
		})
	}
}
