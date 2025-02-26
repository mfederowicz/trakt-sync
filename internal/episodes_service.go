// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// EpisodesService  handles communication with the episodes related
// methods of the Trakt API.
type EpisodesService Service

// GetEpisode Returns episode object.
func (m *EpisodesService) GetEpisode(ctx context.Context, id *int) (*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("episodes/%d", *id)
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

// GetSingleEpisodeForShow Returns a single episode's details.
//
// API docs: https://trakt.docs.apiary.io/#reference/episodes/summary/get-a-single-episode-for-a-show
func (m *EpisodesService) GetSingleEpisodeForShow(ctx context.Context, id *int, season *int, episode *int) (*str.Episode, *str.Response, error) {
	var url = fmt.Sprintf("shows/%d/seasons/%d/episodes/%d", *id, *season, *episode)
	printer.Println("fetch single episode url:" + url)
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
