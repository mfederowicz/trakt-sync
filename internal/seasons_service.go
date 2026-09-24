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

// SeasonsService  handles communication with the seasons related
// methods of the Trakt API.
type SeasonsService Service

// GetSeason Returns season object.
func (s *SeasonsService) GetSeason(ctx context.Context, id *string, opts *uri.ListOptions) (*str.Season, *str.Response, error) {
	var url = fmt.Sprintf("seasons/%s", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch season url:" + url)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, nil, err
	}

	result := new(str.Season)
	resp, err := s.client.Do(ctx, req, &result)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("season not found with traktId:%s", *id)
	}

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ReportSeason Report a season, by its Trakt ID, for moderator review.
//
// API docs: https://docs.trakt.tv/reference/postseasonsreport
func (s *SeasonsService) ReportSeason(ctx context.Context, id *string, report *str.SeasonReport) (*str.Response, error) {
	var url = fmt.Sprintf("seasons/%s/report", *id)
	req, err := s.client.NewRequest(http.MethodPost, url, report)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	var conflict *ConflictError
	if errors.As(err, &conflict) {
		return resp, fmt.Errorf(consts.SeasonIDReportPending, *id)
	}
	if err != nil {
		return resp, err
	}

	return resp, nil
}
