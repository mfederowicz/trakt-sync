// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// RecommendationsService  handles communication with the recommendations related
// methods of the Trakt API.
type RecommendationsService Service

// HideMovieRecommendation to hide a movie from getting recommended anymore..
// API docs:https://trakt.docs.apiary.io/#reference/recommendations/hide-movie/hide-a-movie-recommendation
func (m *RecommendationsService) HideMovieRecommendation(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("recommendations/movies/%s", *id)
	m.client.debug("hide recommendations")
	req, err := m.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(ctx, req, nil)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf(consts.RecommendationNotFoundWithID, *id)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// HideShowRecommendation to hide a show from getting recommended anymore.
// API docs:https://trakt.docs.apiary.io/#reference/recommendations/hide-show/hide-a-show-recommendation
func (m *RecommendationsService) HideShowRecommendation(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("recommendations/shows/%s", *id)
	m.client.debug("hide recommendations")
	req, err := m.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(ctx, req, nil)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf(consts.RecommendationNotFoundWithID, *id)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetMovieRecommendations Movie recommendations for a user.
// API docs:https://trakt.docs.apiary.io/#reference/recommendations/movies/get-movie-recommendations
func (m *RecommendationsService) GetMovieRecommendations(ctx context.Context, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
	var url = "recommendations/movies"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	m.client.debug("fetch recommendations url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Recommendation{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		m.client.debug("fetch recommendations err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetShowRecommendations Show recommendations for a user.
// API docs:https://trakt.docs.apiary.io/#reference/recommendations/shows/get-show-recommendations
func (m *RecommendationsService) GetShowRecommendations(ctx context.Context, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
	var url = "recommendations/shows"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	m.client.debug("fetch recommendations url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Recommendation{}
	resp, err := m.client.Do(ctx, req, &list)

	if err != nil {
		m.client.debug("fetch recommendations err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
