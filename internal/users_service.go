// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
)

// UsersService  handles communication with the users related
// methods of the Trakt API.
type UsersService Service

// GetItemstOnAPersonalList Get all items on a personal list.
//
// API docs: https://trakt.docs.apiary.io/#reference/users/list-items/get-items-on-a-personal-list
func (u *UsersService) GetItemstOnAPersonalList(ctx context.Context, id *string, listID *string, t *string) ([]*str.UserListItem, *str.Response, error) {
	var url string

	if id != nil {
		url = fmt.Sprintf("users/%s/lists/%s/items/%s", *id, *listID, *t)
	} else {
		url = "users/me/lists/watchlist/items/movies"
	}

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.UserListItem{}
	resp, err := u.client.Do(ctx, req, &lists)

	if err != nil {
		fmt.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return lists, resp, nil
}

// GetUsersPersonalLists Returns all personal lists for a user.
//
// API docs: https://trakt.docs.apiary.io/#reference/users/lists/get-a-user's-personal-lists
func (u *UsersService) GetUsersPersonalLists(ctx context.Context, id *string) ([]*str.PersonalList, *str.Response, error) {
	var url string

	if id != nil {
		url = fmt.Sprintf("users/%s/lists", *id)
	} else {
		url = "users/me/lists"
	}

	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.PersonalList{}
	resp, err := u.client.Do(ctx, req, &lists)

	if err != nil {
		fmt.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return lists, resp, nil
}

// GetUserProfile Get a user's profile information.
//
// API docs:https://trakt.docs.apiary.io/#reference/users/profile/get-user-profile
func (u *UsersService) GetUserProfile(ctx context.Context, id *string) (*str.UserProfile, *str.Response, error) {
	var url string

	if id != nil {
		url = fmt.Sprintf("users/%s", *id)
	} else {
		url = "user/me"
	}

	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	profile := new(str.UserProfile)
	resp, err := u.client.Do(ctx, req, &profile)

	if err != nil {
		return nil, resp, err
	}

	return profile, resp, nil
}
