// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// MoviesService  handles communication with the movies related
// methods of the Trakt API.
type MoviesService Service

// GetMovie Returns movie object.
func (m *MoviesService) GetMovie(ctx context.Context, id *int) (*str.Movie, *str.Response, error) {
	var url = fmt.Sprintf("movies/%d", *id)
	printer.Println("fetch movie url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	movie := new(str.Movie)
	resp, err := m.client.Do(ctx, req, &movie)

	if err != nil {
		printer.Println("fetch movie err:" + err.Error())
		return nil, resp, err
	}

	return movie, resp, nil
}
