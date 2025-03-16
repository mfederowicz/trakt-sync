// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// GenresService  handles communication with the genres related
// methods of the Trakt API.
type GenresService Service

// GetGenres Get a list of all genres, including names and slugs.
//
// API docs: https://trakt.docs.apiary.io/#reference/genres/list/get-genres 
func (g *GenresService) GetGenres(ctx context.Context, strType *string) ([]*str.Genre, *str.Response, error) {
	var url = fmt.Sprintf("genres/%s", *strType)
	printer.Println("fetch genres url:" + url)
	
	req, err := g.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	list := []*str.Genre{}
	resp, err := g.client.Do(ctx, req, &list)	

	if err != nil {
		printer.Println("fetch genres err:", err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}


