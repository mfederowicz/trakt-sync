// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
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
	printer.Println("personal list url:" + url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.UserListItem{}
	resp, err := u.client.Do(ctx, req, &lists)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
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

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.PersonalList{}
	resp, err := u.client.Do(ctx, req, &lists)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
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

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
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

// GetSavedFilters Get all saved filters a users has created.
//
// API docs: https://trakt.docs.apiary.io/#reference/users/saved-filters/get-saved-filters
func (u *UsersService) GetSavedFilters(ctx context.Context, section *string) ([]*str.SavedFilter, *str.Response, error) {
	var url string

	url = fmt.Sprintf("users/saved_filters/%s", *section)

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	lists := []*str.SavedFilter{}
	resp, err := u.client.Do(ctx, req, &lists)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return lists, resp, nil
}

// GetStats Returns stats about the movies, shows, and episodes a user has watched, collected, and rated.
//
// API docs:https://trakt.docs.apiary.io/#reference/users/stats/get-stats
func (u *UsersService) GetStats(ctx context.Context, id *string) (*str.UserStats, *str.Response, error) {
	var url string

	if id != nil {
		url = fmt.Sprintf("users/%s/stats", *id)
	} else {
		url = "users/me/stats"
	}

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	stats := new(str.UserStats)
	resp, err := u.client.Do(ctx, req, &stats)

	if err != nil {
		return nil, resp, err
	}

	return stats, resp, nil
}

// GetWatched Returns all movies or shows a user has watched sorted by most plays.
//
// API docs:https://trakt.docs.apiary.io/#reference/users/watched/get-watched
func (u *UsersService) GetWatched(ctx context.Context, id *string, watchType *string, opts *uri.ListOptions) ([]*str.UserWatched, *str.Response, error) {
	var url string

	if id != nil {
		url = fmt.Sprintf("users/%s/watched/%s", *id, *watchType)
	} else {
		url = fmt.Sprintf("users/me/watched/%s", *watchType)
	}

	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	watched := []*str.UserWatched{}
	resp, err := u.client.Do(ctx, req, &watched)

	if err != nil {
		return nil, resp, err
	}

	return watched, resp, nil
}

// RetrieveSettings Get the user's settings so you can align your app's experience with what they're used to on the trakt website.
// API docs: https://trakt.docs.apiary.io/#reference/users/settings/retrieve-settings
func (u *UsersService) RetrieveSettings(ctx context.Context) (*str.UserSettings, *str.Response, error) {
	url := "users/settings"
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	settings := new(str.UserSettings)
	resp, err := u.client.Do(ctx, req, &settings)

	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}
