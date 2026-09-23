// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// TestServicesWithoutResponse checks that services return an error instead of
// panicking when the request fails before any HTTP response (resp is nil).
// One method per service file that checks resp.StatusCode before err.
func TestServicesWithoutResponse(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{}
	tests := []struct {
		name string
		call func(c *Client) error
	}{
		{name: "comments", call: func(c *Client) error { _, _, err := c.Comments.PostAComment(ctx, &str.Comment{}); return err }},
		{name: "movies", call: func(c *Client) error { _, _, err := c.Movies.GetMovie(ctx, str.String("tron"), opts); return err }},
		{name: "networks", call: func(c *Client) error { _, _, err := c.Networks.GetNetworksList(ctx, opts); return err }},
		{name: "notes", call: func(c *Client) error { _, _, err := c.Notes.AddNotes(ctx, &str.Notes{}); return err }},
		{name: "people", call: func(c *Client) error {
			_, _, err := c.People.GetAllPeopleForShow(ctx, str.String("bb"), opts)
			return err
		}},
		{name: "recommendations", call: func(c *Client) error {
			_, err := c.Recommendations.HideMovieRecommendation(ctx, str.String("tron"))
			return err
		}},
		{name: "seasons", call: func(c *Client) error { _, _, err := c.Seasons.GetSeason(ctx, str.String("bb"), opts); return err }},
		{name: "shows", call: func(c *Client) error { _, _, err := c.Shows.GetShow(ctx, str.String("bb"), opts); return err }},
		{name: "sync", call: func(c *Client) error { _, err := c.Sync.RemovePlaybackItem(ctx, new(int)); return err }},
		{name: "users", call: func(c *Client) error {
			_, _, err := c.Users.AddPersonalList(ctx, str.String("sean"), &str.PersonalList{})
			return err
		}},
		{name: "users unfollow", call: func(c *Client) error { _, err := c.Users.Unfollow(ctx, str.String("sean")); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			setup.Teardown() // closed server: every request fails without a response

			if err := tt.call(setup.Client); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

// TestPostWithoutValidationErrors checks that a non-422 error on POST comments/notes
// returns the HTTP error instead of panicking on the missing validation errors.
func TestPostWithoutValidationErrors(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		path string
		call func(c *Client) error
	}{
		{name: "comments", path: "/comments", call: func(c *Client) error { _, _, err := c.Comments.PostAComment(ctx, &str.Comment{}); return err }},
		{name: "notes", path: "/notes", call: func(c *Client) error { _, _, err := c.Notes.AddNotes(ctx, &str.Notes{}); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			})

			err := tt.call(setup.Client)
			var invalidUser *InvalidUserError
			if !errors.As(err, &invalidUser) {
				t.Fatalf("error is %v, want *InvalidUserError", err)
			}
		})
	}
}
