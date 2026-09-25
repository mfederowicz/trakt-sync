// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SmartListsService handles communication with the smart lists related
// methods of the Trakt API.
type SmartListsService Service

// GetSmartList Returns a single smart list definition by its slug.
//
// API docs: https://docs.trakt.tv/reference/getsmart_listssummary
func (s *SmartListsService) GetSmartList(ctx context.Context, id *string) (*str.SmartList, *str.Response, error) {
	var url = fmt.Sprintf("smart-lists/%s", *id)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.SmartList)
	resp, err := s.client.Do(ctx, req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetSmartListItems Returns the dynamic items a smart list resolves to.
//
// API docs: https://docs.trakt.tv/reference/getsmart_listsitems
func (s *SmartListsService) GetSmartListItems(ctx context.Context, id *string, opts *uri.SmartListItemsOptions) ([]*str.UserListItem, *str.Response, error) {
	var url = fmt.Sprintf("smart-lists/%s/items", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.UserListItem{}
	resp, err := s.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
