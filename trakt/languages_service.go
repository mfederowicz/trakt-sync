// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
)

// LanguagesService  handles communication with the languages related
// methods of the Trakt API.
type LanguagesService Service

// GetLanguages Get a list of all languages, including names and codes.
//
// API docs: https://trakt.docs.apiary.io/#reference/languages/list/get-languages
func (g *LanguagesService) GetLanguages(ctx context.Context, strType *string) ([]*str.Language, *str.Response, error) {
	var url = fmt.Sprintf("languages/%s", *strType)
	g.client.debug("fetch languages url:" + url)

	req, err := g.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	list := []*str.Language{}
	resp, err := g.client.Do(ctx, req, &list)

	if err != nil {
		g.client.debug("fetch languages err:", err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}
