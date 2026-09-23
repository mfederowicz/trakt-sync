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

// ListsService  handles communication with the lists related
// methods of the Trakt API.
type ListsService Service

// GetTrendingLists Returns all lists with the most likes and comments over the last 7 days.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/trending/get-trending-lists
func (l *ListsService) GetTrendingLists(ctx context.Context, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	var url string

	url = "lists/trending"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch trending url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.List{}
	resp, err := l.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetPopularLists Returns the most popular lists. Popularity is calculated using total number of likes and comments..
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/popular/get-popular-lists
func (l *ListsService) GetPopularLists(ctx context.Context, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	var url string

	url = "lists/popular"
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch trending url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.List{}
	resp, err := l.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetList Returns a single list. Use the /lists/:id/items method to get the actual items this list contains.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/list/get-list
func (l *ListsService) GetList(ctx context.Context, id *string) (*str.PersonalList, *str.Response, error) {
	var url = fmt.Sprintf("lists/%s", *id)
	printer.Println("fetch single list:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := new(str.PersonalList)
	resp, err := l.client.Do(ctx, req, &list)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("list not found with traktId:%s", *id)
	}

	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}

// GetAllUsersWhoLikedList Returns all users who liked a list.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/list-likes/get-all-users-who-liked-a-list
func (l *ListsService) GetAllUsersWhoLikedList(ctx context.Context, opts *uri.ListOptions, id *string) ([]*str.UserLike, *str.Response, error) {
	var url = fmt.Sprintf("lists/%s/likes", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch likes url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.UserLike{}
	resp, err := l.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// LikeList Votes help determine popular lists. Only one like is allowed per list per user.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/list-like/like-a-list
func (l *ListsService) LikeList(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("lists/%s/like", *id)
	printer.Println("send like for single list:" + url)
	req, err := l.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveLikeList Remove a like on a list.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/list-like/remove-like-on-a-list
func (l *ListsService) RemoveLikeList(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("lists/%s/like", *id)
	printer.Println("remove like for single list:" + url)
	req, err := l.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetListItems Returns items from single list.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/list-items/get-items-on-a-list
func (l *ListsService) GetListItems(ctx context.Context, id *string, t *string, opts *uri.ListOptions) ([]*str.UserListItem, *str.Response, error) {
	var url string

	if t != nil {
		url = fmt.Sprintf("lists/%s/items/%s", *id, *t)
	} else {
		url = fmt.Sprintf("lists/%s/items", *id)
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("list url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.UserListItem{}
	resp, err := l.client.Do(ctx, req, &lists)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return lists, resp, nil
}

// GetListComments Returns comments from single list.
//
// API docs: https://trakt.docs.apiary.io/#reference/lists/list-comments/get-all-list-comments
func (l *ListsService) GetListComments(ctx context.Context, id *string, sort *string, opts *uri.ListOptions) ([]*str.ListComment, *str.Response, error) {
	var url string

	if sort != nil {
		url = fmt.Sprintf("lists/%s/comments/%s", *id, *sort)
	} else {
		url = fmt.Sprintf("lists/%s/comments", *id)
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("list url:" + url)
	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.ListComment{}
	resp, err := l.client.Do(ctx, req, &lists)

	if err != nil {
		// if it's just a 404, don't return that as an error
		_, err = parseBoolResponse(err)
		return nil, resp, err
	}

	return lists, resp, nil
}

// GetTrendingListsByType Returns trending lists of the given type, ordered by current activity.
//
// API docs: https://docs.trakt.tv/reference/getliststrendingbytype
func (l *ListsService) GetTrendingListsByType(ctx context.Context, listType *string, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	return l.fetchLists(ctx, fmt.Sprintf("lists/trending/%s", *listType), opts)
}

// GetPopularListsByType Returns popular lists of the given type, ordered by long-term activity.
//
// API docs: https://docs.trakt.tv/reference/getlistspopularbytype
func (l *ListsService) GetPopularListsByType(ctx context.Context, listType *string, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	return l.fetchLists(ctx, fmt.Sprintf("lists/popular/%s", *listType), opts)
}

// ReportList Report a list for moderator review.
//
// API docs: https://docs.trakt.tv/reference/postlistsreport
func (l *ListsService) ReportList(ctx context.Context, id *string, report *str.ListReport) (*str.Response, error) {
	var url = fmt.Sprintf("lists/%s/report", *id)
	req, err := l.client.NewRequest(http.MethodPost, url, report)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(ctx, req, nil)
	var conflict *ConflictError
	if errors.As(err, &conflict) {
		return resp, fmt.Errorf(consts.ListReportPending, *id)
	}
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (l *ListsService) fetchLists(ctx context.Context, url string, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := l.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.List{}
	resp, err := l.client.Do(ctx, req, &list)
	if err != nil {
		return nil, resp, err
	}

	return list, resp, nil
}
