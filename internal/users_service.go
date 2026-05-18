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

// GetPendingFollowingRequests List a user's pending following requests that they're waiting for the other user's to approve.
// API docs:https://trakt.docs.apiary.io/#reference/users/following-requests/get-pending-following-requests
func (u *UsersService) GetPendingFollowingRequests(ctx context.Context, options *uri.ListOptions) ([]*str.FollowRequest, *str.Response, error) {
	url := "users/requests/following"
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	url, err = uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}

	requests := []*str.FollowRequest{}
	resp, err := u.client.Do(ctx, req, &requests)

	if err != nil {
		return nil, resp, err
	}

	return requests, resp, nil
}

// GetFollowRequests List a user's pending follow requests so they can either approve or deny them.
// API docs:https://trakt.docs.apiary.io/#reference/users/follower-requests/get-follow-requests
func (u *UsersService) GetFollowRequests(ctx context.Context, options *uri.ListOptions) ([]*str.FollowRequest, *str.Response, error) {
	url := "users/requests/following"
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	url, err = uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}

	requests := []*str.FollowRequest{}
	resp, err := u.client.Do(ctx, req, &requests)

	if err != nil {
		return nil, resp, err
	}

	return requests, resp, nil
}

// ApproveFollowRequest Approve a follower using the id of the request.
// If the id is not found, was already approved, or was already denied, a 404 error will be returned.
// API docs:https://trakt.docs.apiary.io/#reference/users/approve-or-deny-follower-requests/approve-follow-request
func (u *UsersService) ApproveFollowRequest(ctx context.Context, request int) (*str.FollowRequest, *str.Response, error) {
	var url string

	url = fmt.Sprintf("users/requests/%d", *&request)

	printer.Println("approve follower")
	req, err := u.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, nil, err
	}

	fr := new(str.FollowRequest)
	resp, err := u.client.Do(ctx, req, fr)
	if err != nil {
		return fr, resp, err
	}

	return fr, resp, nil
}

// DenyFollowRequest Deny a follower using the id of the request.
// If the id is not found, was already approved, or was already denied, a 404 error will be returned.
func (u *UsersService) DenyFollowRequest(ctx context.Context, request int) (*str.FollowRequest, *str.Response, error) {
	var url string

	url = fmt.Sprintf("users/requests/%d", *&request)

	printer.Println("deny follower")
	req, err := u.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, nil, err
	}

	fr := new(str.FollowRequest)
	resp, err := u.client.Do(ctx, req, fr)
	if err != nil {
		return fr, resp, err
	}

	return fr, resp, nil
}

// GetHiddenItems Get hidden items for a section. This will return an array of
// standard media objects. You can optionally limit the type of results to return..
// API docs:https:https://trakt.docs.apiary.io/#reference/users/hidden-items/get-hidden-items
func (u *UsersService) GetHiddenItems(ctx context.Context, section *string, opts *uri.ListOptions) ([]*str.HiddenItem, *str.Response, error) {
	var url string

	if section != nil {
		url = fmt.Sprintf("users/hidden/%s", *section)
	} else {
		url = "users/hidden"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.HiddenItem{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// AddHiddenItems Hide items for a specific section. Here's what type of items
// can hidden for each section. You can optionally specify the
// hidden_at date for each item.
// API docs:https://trakt.docs.apiary.io/#reference/users/add-hidden-items/add-hidden-items
func (u *UsersService) AddHiddenItems(ctx context.Context, items *str.HistoryItems, section string) (*str.AddResult, error) {
	var url string
	url = fmt.Sprintf("users/hidden/%s", section)
	printer.Println("add hidden items")
	req, err := u.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.AddResult)
	_, err = u.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// RemoveHiddenItems Unhide items for a specific section. Here's what type of items can unhidden for each section.
// API docs:https://trakt.docs.apiary.io/#reference/users/remove-hidden-items/remove-hidden-items
func (u *UsersService) RemoveHiddenItems(ctx context.Context, items *str.HistoryItems, section string) (*str.RemoveResult, error) {
	var url string
	url = fmt.Sprintf("users/hidden/%s/remove", section)
	printer.Println("remove hidden items")
	req, err := u.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, err
	}

	result := new(str.RemoveResult)
	_, err = u.client.Do(ctx, req, result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetProfile Get a user's profile information. If the user is private,
// info will only be returned if you send OAuth and are either that user
// or an approved follower. Adding ?extended=vip will return some additional VIP related fields
// so you can display the user's Trakt VIP status and year count.
// API docs:https://trakt.docs.apiary.io/#reference/users/profile/get-user-profile
func (u *UsersService) GetProfile(ctx context.Context, s *string) (*str.UserProfile, *str.Response, error) {
	url := fmt.Sprintf("users/%s", *s)
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

// GetLikes Get items a user likes. This will return an array of standard media objects.
// You can optionally limit the type of results to return.
// API docs:https://trakt.docs.apiary.io/#reference/users/likes/get-likes
func (u *UsersService) GetLikes(ctx context.Context, user *string, stype *string, opts *uri.ListOptions) ([]*str.UserLike, *str.Response, error) {
	var url string
	if stype != nil {
		url = fmt.Sprintf("users/%s/likes/%s", *user, *stype)
	} else {
		url = "users/me/likes"
	}
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.UserLike{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}
