// Package internal used for client and services
package internal

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// TeamService handles communication with the team related
// methods of the Trakt API.
type TeamService Service

// GetTeamMembers Returns Trakt team members.
//
// API docs: https://docs.trakt.tv/reference/getteammembers
func (t *TeamService) GetTeamMembers(ctx context.Context, opts *uri.ListOptions) ([]*str.TeamMember, *str.Response, error) {
	var url = "team"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := t.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.TeamMember{}
	resp, err := t.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
