package trakt

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SocialRecommendationsService handles communication with the social recommendations related
// methods of the Trakt API.
type SocialRecommendationsService Service

// GetSocialMovieRecommendations Returns movie recommendations based on the authenticated user social graph.
//
// API docs: https://docs.trakt.tv/reference/getsocial_recommendationsmoviesrecommend
//
// Deprecated: Trakt dropped this endpoint from its API contract on 2026-10-13. The method still sends its
// request; in a later release it will return an *EndpointRemovedError instead, and then be removed.
func (s *SocialRecommendationsService) GetSocialMovieRecommendations(ctx context.Context, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
	if err := lockedEndpoint(http.MethodGet, "social_recommendations/movies"); err != nil {
		return nil, nil, err
	}
	return s.fetchRecommendations(ctx, "social_recommendations/movies", opts)
}

// GetSocialShowRecommendations Returns show recommendations based on the authenticated user social graph.
//
// API docs: https://docs.trakt.tv/reference/getsocial_recommendationsshowsrecommend
//
// Deprecated: Trakt dropped this endpoint from its API contract on 2026-10-13. The method still sends its
// request; in a later release it will return an *EndpointRemovedError instead, and then be removed.
func (s *SocialRecommendationsService) GetSocialShowRecommendations(ctx context.Context, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
	if err := lockedEndpoint(http.MethodGet, "social_recommendations/shows"); err != nil {
		return nil, nil, err
	}
	return s.fetchRecommendations(ctx, "social_recommendations/shows", opts)
}

func (s *SocialRecommendationsService) fetchRecommendations(ctx context.Context, url string, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Recommendation{}
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
