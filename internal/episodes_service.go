// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// EpisodesService  handles communication with the episodes related
// methods of the Trakt API.
type EpisodesService Service

// GetEpisode Returns episode object.
func (m *EpisodesService) GetEpisode(ctx context.Context, id *string) (*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("episodes/%s", *id)
	printer.Println("fetch episode url:" + url)
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Episode)
	resp, err := m.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch episode err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// ReportEpisode Report an episode, by its Trakt ID, for moderator review.
//
// API docs: https://docs.trakt.tv/reference/postepisodesreport
func (m *EpisodesService) ReportEpisode(ctx context.Context, id *string, report *str.EpisodeReport) (*str.Response, error) {
	var url = fmt.Sprintf("episodes/%s/report", *id)
	req, err := m.client.NewRequest(http.MethodPost, url, report)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(ctx, req, nil)
	var conflict *ConflictError
	if errors.As(err, &conflict) {
		return resp, fmt.Errorf(consts.EpisodeIDReportPending, *id)
	}
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetEpisodeWatchNow Returns streaming and watch now sources for an episode, by its Trakt ID, in the requested country.
//
// API docs: https://docs.trakt.tv/reference/getepisodeswatchnow
func (m *EpisodesService) GetEpisodeWatchNow(ctx context.Context, id *string, country *string, opts *uri.ListOptions) (map[string]*str.WatchNowSources, *str.Response, error) {
	var url = fmt.Sprintf("episodes/%s/watchnow/%s", *id, *country)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := m.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := map[string]*str.WatchNowSources{}
	resp, err := m.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
