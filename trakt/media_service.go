package trakt

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// MediaService handles communication with the media related
// methods of the Trakt API.
type MediaService Service

// GetTrendingMedia Returns trending movies and shows, ordered by current watcher activity.
//
// API docs: https://docs.trakt.tv/reference/getmediatrending
func (m *MediaService) GetTrendingMedia(ctx context.Context, opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error) {
	list := []*str.MediaItem{}
	resp, err := m.fetchList(ctx, "media/trending", opts, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPopularMedia Returns popular movies and shows.
//
// API docs: https://docs.trakt.tv/reference/getmediapopular
func (m *MediaService) GetPopularMedia(ctx context.Context, opts *uri.ListOptions) ([]*str.Media, *str.Response, error) {
	list := []*str.Media{}
	resp, err := m.fetchList(ctx, "media/popular", opts, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAnticipatedMedia Returns anticipated movies and shows, based on list activity.
//
// API docs: https://docs.trakt.tv/reference/getmediaanticipated
func (m *MediaService) GetAnticipatedMedia(ctx context.Context, opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error) {
	list := []*str.MediaItem{}
	resp, err := m.fetchList(ctx, "media/anticipated", opts, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAnticipatedRecommendations Returns unreleased movies and shows ranked for the taste of the user, at most 100 in total.
//
// API docs: https://docs.trakt.tv/reference/getmediaanticipatedrecommendations
func (m *MediaService) GetAnticipatedRecommendations(ctx context.Context, opts *uri.ListOptions) ([]*str.AnticipatedRecommendation, *str.Response, error) {
	var url = "media/recommendations/anticipated"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.AnticipatedRecommendation{}
	resp, err := m.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

func (m *MediaService) fetchList(ctx context.Context, url string, opts *uri.ListOptions, list any) (*str.Response, error) {
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return m.client.Do(ctx, req, list)
}
