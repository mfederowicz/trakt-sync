// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

// YounifyService handles communication with the younify (streaming service connections) related
// methods of the Trakt API.
type YounifyService Service

// GetConnections Lists every connectable streaming service with the user's connection status.
//
// API docs: https://docs.trakt.tv/reference/getyounifyconnections
func (y *YounifyService) GetConnections(ctx context.Context) ([]*str.YounifyConnection, *str.Response, error) {
	req, err := y.client.NewRequest(http.MethodGet, "younify/connections", nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.YounifyConnection{}
	resp, err := y.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// Connect Creates a signed web auth URL to connect a streaming service.
//
// API docs: https://docs.trakt.tv/reference/postyounifyconnect
func (y *YounifyService) Connect(ctx context.Context, connect *str.YounifyConnect) (*str.YounifyConnectResult, *str.Response, error) {
	req, err := y.client.NewRequest(http.MethodPost, "younify/connect", connect)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.YounifyConnectResult)
	resp, err := y.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RefreshService Queues a re-sync of a connected streaming service; allData forces a full re-sync.
//
// API docs: https://docs.trakt.tv/reference/postyounifyrefresh
// API docs: https://docs.trakt.tv/reference/postyounifyrefreshall
func (y *YounifyService) RefreshService(ctx context.Context, serviceID *string, allData bool) (*str.Response, error) {
	var url = fmt.Sprintf("younify/users/refresh/%s", *serviceID)
	if allData {
		url = fmt.Sprintf("%s/%s", url, consts.AllDataSegment)
	}
	req, err := y.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	return y.client.Do(ctx, req, nil)
}

// DisconnectService Unlinks a streaming service from the user.
//
// API docs: https://docs.trakt.tv/reference/deleteyounifydisconnect
func (y *YounifyService) DisconnectService(ctx context.Context, serviceID *string) (*str.Response, error) {
	var url = fmt.Sprintf("younify/users/services/%s", *serviceID)
	req, err := y.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return y.client.Do(ctx, req, nil)
}
