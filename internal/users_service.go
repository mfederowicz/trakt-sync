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

	fmt.Println("url:", url)
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

// GetListLikes Returns all users who liked a list.
// API docs:https://trakt.docs.apiary.io/#reference/users/list-likes/get-all-users-who-liked-a-list
func (u *UsersService) GetListLikes(ctx context.Context, user *string, listID *string, opts *uri.ListOptions) ([]*str.UserLike, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/likes", *user, *listID)
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

// GetCollection Get all collected items in a user's collection.
// A collected item indicates availability to watch digitally or on physical media.
// API docs:https://trakt.docs.apiary.io/#reference/users/collection/get-collection
func (u *UsersService) GetCollection(ctx context.Context, user *string, stype *string, opts *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string
	if stype != nil {
		url = fmt.Sprintf("users/%s/collection/%s", *user, *stype)
	} else {
		url = "users/me/collection"
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

	items := []*str.ExportlistItem{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetComments Returns the most recently written comments for the user.
// You can optionally filter by the comment_type and media type to limit what gets returned.
// By default, only top level comments are returned. Set ?include_replies=true to return
// replies in addition to top level comments. Set ?include_replies=only to return only
// replies and no top level comments.
// API docs:https://trakt.docs.apiary.io/#reference/users/comments/get-comments
func (u *UsersService) GetComments(ctx context.Context, user *string, commentType *string, strType *string, opts *uri.ListOptions) ([]*str.CommentItem, *str.Response, error) {
	var url string
	if commentType != nil && strType != nil {
		url = fmt.Sprintf("users/%s/comments/%s/%s", *user, *commentType, *strType)
	} else {
		url = "users/me/comments/all/all"
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

	items := []*str.CommentItem{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetNotes Returns the most recently notes for the user.
// You can optionally filter by media type to limit what gets returned.
// Use the attached_to info to know what the note is actually added to.
// Media items like movie, show, season, episode, or person are straightforward,
// but history will need to be mapped to that specific play in their watched history
// since they might have multiple plays. Since collection and rating is a 1:1 association,
// you can assume the note is attached to the media item in the type field that has been collected or rated.
// API docs:https://trakt.docs.apiary.io/#reference/users/notes/get-notes
func (u *UsersService) GetNotes(ctx context.Context, user *string, strType *string, opts *uri.ListOptions) ([]*str.NotesItem, *str.Response, error) {
	var url string
	if strType != nil {
		url = fmt.Sprintf("users/%s/notes/%s", *user, *strType)
	} else {
		url = "users/me/notes/all"
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

	items := []*str.NotesItem{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// AddPersonalList Create a new personal list. The name is the only required field, but the other info is recommended to ask for.
// API docs:https://trakt.docs.apiary.io/#reference/users/lists/create-personal-list
func (u *UsersService) AddPersonalList(ctx context.Context, user *string, list *str.PersonalList) (*str.PersonalList, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists", *user)
	printer.Println("create new personal list")
	req, err := u.client.NewRequest(http.MethodPost, url, list)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.PersonalList)
	resp, err := u.client.Do(ctx, req, result)

	if resp.StatusCode == 420 {
		return nil, nil, errors.New("use the /users/settings method to get all limits for a user account. In most cases, upgrading to Trakt VIP will increase the limits")
	}

	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// ReorderLists Reorder all lists by sending the updated rank of list ids. Use the /users/:id/lists method to get all list ids.
// API docs:https://trakt.docs.apiary.io/#reference/users/reorder-lists/reorder-a-user's-lists
func (u *UsersService) ReorderLists(ctx context.Context, user *string, items *str.ItemsToReorder) (*str.ReorderResults, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/reorder", *user)
	printer.Println("reorder user lists")
	req, err := u.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.ReorderResults)
	resp, err := u.client.Do(ctx, req, result)

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetCollaborations Returns all lists a user can collaborate on.
// This gives full access to add, remove, and re-order list items.
// It essentially works just like a list owned by the user, just make sure to
// use the correct list owner user when building the API URLs.
// API docs:https://trakt.docs.apiary.io/#reference/users/collaborations/get-all-lists-a-user-can-collaborate-on
func (u *UsersService) GetCollaborations(ctx context.Context, user *string, opts *uri.ListOptions) ([]*str.PersonalList, *str.Response, error) {
	var url string
	if user != nil {
		url = fmt.Sprintf("users/%s/lists/collaborations", *user)
	} else {
		url = "users/me/lists/collaborations"
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

	items := []*str.PersonalList{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetList Returns a single personal list. Use the /users/:id/lists/:list_id/items method to get the actual items this list contains.
// API docs:https://trakt.docs.apiary.io/#reference/users/list/get-personal-list
func (u *UsersService) GetList(ctx context.Context, user *string, listID *string, opts *uri.ListOptions) (*str.PersonalList, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s", *user, *listID)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	// fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	item := new(str.PersonalList)
	resp, err := u.client.Do(ctx, req, &item)

	if err != nil {
		return nil, resp, err
	}

	return item, resp, nil
}

// UpdateList Update a personal list by sending 1 or more parameters.
// If you update the list name, the original slug will still be retained
// so existing references to this list won't break.
// API docs:https://trakt.docs.apiary.io/#reference/users/list/update-personal-list
func (u *UsersService) UpdateList(ctx context.Context, user *string, listID *string, update *str.PersonalList) (*str.PersonalList, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s", *user, *listID)
	req, err := u.client.NewRequest(http.MethodPut, url, update)
	if err != nil {
		return nil, nil, err
	}

	item := new(str.PersonalList)
	resp, err := u.client.Do(ctx, req, &item)

	if err != nil {
		return nil, resp, err
	}

	return item, resp, nil
}

// DeleteList Remove a personal list and all items it contains.
// API docs:https://trakt.docs.apiary.io/#reference/users/list/delete-a-user's-personal-list
func (u *UsersService) DeleteList(ctx context.Context, user *string, listID *string) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s", *user, *listID)
	req, err := u.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveListLike Remove a like on a list.
// API docs:https://trakt.docs.apiary.io/#reference/users/list-like/remove-like-on-a-list
func (u *UsersService) RemoveListLike(ctx context.Context, user *string, listID *string) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/like", *user, *listID)
	req, err := u.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ListLike Votes help determine popular lists. Only one like is allowed per list per user.
// API docs:https://trakt.docs.apiary.io/#reference/users/list-like/like-a-list
func (u *UsersService) ListLike(ctx context.Context, user *string, listID *string) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/like", *user, *listID)
	req, err := u.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetListItems Get all items on a personal list. Items can be a movie, show, season, episode, or person.
// You can optionally specify the type parameter with a single value or comma delimited string for multiple item types.
// API docs:https://trakt.docs.apiary.io/#reference/users/list-items/get-items-on-a-personal-list
func (u *UsersService) GetListItems(ctx context.Context, user *string, listID *string, strType *string, sortBy *string, sortHow *string, options *uri.ListOptions) ([]*str.UserListItem, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/items/%s/%s/%s", *user, *listID, *strType, *sortBy, *sortHow)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.UserListItem{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// AddListItems Add one or more items to a personal list. Items can be movies, shows, seasons, episodes, or people.
// API docs:https://trakt.docs.apiary.io/#reference/users/add-list-items/add-items-to-personal-list
func (u *UsersService) AddListItems(ctx context.Context, user *string, listID *string, items *str.HistoryItems) (*str.AddResult, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/items", *user, *listID)
	req, err := u.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.AddResult)
	resp, err := u.client.Do(ctx, req, &result)

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RemoveListItems Remove one or more items from a personal list.
// API docs:https://trakt.docs.apiary.io/#reference/users/remove-list-items/remove-items-from-personal-list
func (u *UsersService) RemoveListItems(ctx context.Context, user *string, listID *string, items *str.HistoryItems) (*str.RemoveResult, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/items/remove", *user, *listID)
	req, err := u.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.RemoveResult)
	resp, err := u.client.Do(ctx, req, &result)

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ReorderListItems Reorder all items on a list by sending the updated rank of list item ids.
// Use the /users/:id/lists/:list_id/items method to get all list item ids.
// API docs:https://trakt.docs.apiary.io/#reference/users/reorder-list-items/reorder-items-on-a-list
func (u *UsersService) ReorderListItems(ctx context.Context, user *string, listID *string, items *str.ItemsToReorder) (*str.ReorderResults, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/items/reorder", *user, *listID)
	printer.Println("reorder list items")
	req, err := u.client.NewRequest(http.MethodPost, url, items)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.ReorderResults)
	resp, err := u.client.Do(ctx, req, result)

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateListItem Update the notes on a single list item.
// API docs:https://trakt.docs.apiary.io/#reference/users/update-list-item/update-a-list-item
func (u *UsersService) UpdateListItem(ctx context.Context, user *string, listID *string, listItemID *int, item *str.PersonalListItem) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/items/%d", *user, *listID, *listItemID)
	printer.Println("update list item")
	req, err := u.client.NewRequest(http.MethodPut, url, item)
	if err != nil {
		return nil, err
	}

	result := new(str.PersonalListItem)
	resp, err := u.client.Do(ctx, req, result)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetListComments Returns all top level comments for a list.
// By default, the comments are sorted by most likes.
// Other sorting options include likes_30, most replies, replies_30,
// most plays, highest rating, and added date.
// API docs:https://trakt.docs.apiary.io/#reference/users/list-comments/get-all-list-comments
func (u *UsersService) GetListComments(ctx context.Context, user *string, listID *string, sort *string, options *uri.ListOptions) ([]*str.ListComment, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/comments/%s", *user, *listID, *sort)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.ListComment{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// ListReport Report a user's list for moderator review.
// Send a reason and optional message with additional context.
// A user can only have one pending report per list.
// API docs:https://trakt.docs.apiary.io/#reference/users/list-report/report-a-user's-list
func (u *UsersService) ListReport(ctx context.Context, user *string, listID *string, report *str.ListReport) (*str.ListReportResult, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/lists/%s/report", *user, *listID)
	printer.Println("list report")
	req, err := u.client.NewRequest(http.MethodPost, url, report)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.ListReportResult)
	resp, err := u.client.Do(ctx, req, &result)

	if err != nil {
		return nil, resp, errors.New(*result.Message)
	}

	return result, resp, nil
}

// Follow If the user has a private profile, the follow request will require approval (approved_at will be null).
// If a user is public, they will be followed immediately (approved_at will have a date).
// API docs:https://trakt.docs.apiary.io/#reference/users/follow/follow-this-user
func (u *UsersService) Follow(ctx context.Context, user *string) (*str.FollowResult, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/follow", *user)
	printer.Println("follow user")
	req, err := u.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.FollowResult)
	resp, err := u.client.Do(ctx, req, &result)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

// Unfollow Unfollow someone you already follow..
// API docs:https://trakt.docs.apiary.io/#reference/users/follow/unfollow-this-user
func (u *UsersService) Unfollow(ctx context.Context, user *string) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/follow", *user)
	printer.Println("unfollow user")
	req, err := u.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.client.Do(ctx, req, nil)

	if resp.StatusCode == http.StatusNotFound {
		return resp, nil
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetBlockedUsers Returns all users you have blocked, including when each user was blocked.
// API docs:https://trakt.docs.apiary.io/#reference/users/blocked-users/get-blocked-users
func (u *UsersService) GetBlockedUsers(ctx context.Context, options *uri.ListOptions) ([]*str.UserBlocked, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/blocked")
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.UserBlocked{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// Block Block a user. If they are already following you, they will be removed from your followers.
// Any pending follow request from this user will be blocked, preventing them
// from following you in the future until you unblock them.
// API docs:https://trakt.docs.apiary.io/#reference/users/block/block-this-user
func (u *UsersService) Block(ctx context.Context, user *string) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/block", *user)
	printer.Println("block user")
	req, err := u.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(ctx, req, nil)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Unblock Unblock a user you previously blocked.
// API docs:https://trakt.docs.apiary.io/#reference/users/block/unblock-this-user
func (u *UsersService) Unblock(ctx context.Context, user *string) (*str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/block", *user)
	printer.Println("unblock user")
	req, err := u.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.client.Do(ctx, req, nil)

	if resp.StatusCode == http.StatusNotFound {
		return resp, nil
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetFollowers Returns all followers including when the relationship began.
// API docs:https://trakt.docs.apiary.io/#reference/users/followers/get-followers
func (u *UsersService) GetFollowers(ctx context.Context, user *string, options *uri.ListOptions) ([]*str.Follower, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/followers", *user)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.Follower{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetFollowing Returns all user's they follow including when the relationship began.
// API docs:https://trakt.docs.apiary.io/#reference/users/following/get-following
func (u *UsersService) GetFollowing(ctx context.Context, user *string, options *uri.ListOptions) ([]*str.Follower, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/following", *user)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.Follower{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetFriends Returns all friends for a user including when the relationship began.
// Friendship is a 2 way relationship where each user follows the other.
// API docs:https://trakt.docs.apiary.io/#reference/users/friends/get-friends
func (u *UsersService) GetFriends(ctx context.Context, user *string, options *uri.ListOptions) ([]*str.Friend, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/friends", *user)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.Friend{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetHistory Returns movies and episodes that a user has watched, sorted by most recent.
// You can optionally limit the type to movies or episodes. The id (64-bit integer) in each
// history item uniquely identifies the event and can be used to remove individual events
// by using the /sync/history/remove method. The action will be set to scrobble, checkin,
// or watch.Specify a type and trakt item_id to limit the history for just that item.
// If the item_id is valid, but there is no history, an empty array will be returned.
// API docs:https://trakt.docs.apiary.io/#reference/users/history/get-watched-history
func (u *UsersService) GetHistory(ctx context.Context, user *string, strType *string, id *int, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	if *id > consts.ZeroValue {
		url = fmt.Sprintf("users/%s/history/%s/%d", *user, *strType, *id)
	} else {
		url = fmt.Sprintf("users/%s/history/%s", *user, *strType)
	}

	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(url)
	req, err := u.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	items := []*str.ExportlistItem{}
	resp, err := u.client.Do(ctx, req, &items)

	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

// GetRatings Get a user's ratings filtered by type. You can optionally filter
// for a specific rating between 1 and 10. Send a comma separated string for
// rating if you need multiple ratings.
// API docs:https://trakt.docs.apiary.io/#reference/users/ratings/get-ratings
func (u *UsersService) GetRatings(ctx context.Context, user *string, strType *string, rating *string, options *uri.ListOptions) ([]*str.RatingListItem, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/ratings/%s", *user, *strType)
	if len(*rating) > consts.ZeroValue {
		url = fmt.Sprintf("users/%s/ratings/%s/%s", *user, *strType, *rating)
	}

	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch ratings url:" + url)
	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.RatingListItem{}
	resp, err := u.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchlist Returns all items in a user's watchlist filtered by type.
// API docs:https://trakt.docs.apiary.io/#reference/users/watchlist/get-watchlist
func (u *UsersService) GetWatchlist(ctx context.Context, user *string, types *string, sortBy *string, sortHow *string, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("users/%s/watchlist/%s/%s/%s", *user, *types, *sortBy, *sortHow)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch watchlist url:" + url)
	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := u.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch lists err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetWatchlistComments Returns all top level comments for the watchlist.
// By default, the comments are sorted by most likes.
// Other sorting options include likes_30, most replies, replies_30, most plays, highest rating, and added date.
// API docs:https://trakt.docs.apiary.io/#reference/users/watchlist-comments/get-all-favorites-comments
func (u *UsersService) GetWatchlistComments(ctx context.Context, user *string, sort *string, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string

	url = fmt.Sprintf("users/%s/watchlist/comments/%s", *user, *sort)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch watchlist comments url:" + url)
	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := u.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch watchlist comments err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetFavorites Returns the top 100 shows and movies a user has favorited.
// Apps should encourage user's to add favorites so the algorithm keeps getting better.
// API docs:https://trakt.docs.apiary.io/#reference/users/favorites/get-favorites
func (u *UsersService) GetFavorites(ctx context.Context, user *string, strType *string, sortBy *string, sortHow *string, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/favorites/%s/%s/%s", *user, *strType, *sortBy, *sortHow)

	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch favorites url:" + url)
	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := u.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch favorites err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetFavoritesComments Returns all top level comments for the favorites.
// By default, the comments are sorted by most likes. Other sorting options include
// likes_30, most replies, replies_30, most plays, highest rating, and added date.
// API docs:https://trakt.docs.apiary.io/#reference/users/favorites-comments/get-all-favorites-comments
func (u *UsersService) GetFavoritesComments(ctx context.Context, user *string, sort *string, options *uri.ListOptions) ([]*str.ExportlistItem, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/favorites/comments/%s", *user, *sort)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch favorites comments url:" + url)
	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.ExportlistItem{}
	resp, err := u.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch favorites comments err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// Watching Returns a movie or episode if the user is currently watching something.
// If they are not, it returns no data and a 204 HTTP status code.
// API docs:https://trakt.docs.apiary.io/#reference/users/watching/get-watching
func (u *UsersService) Watching(ctx context.Context, user *string, options *uri.ListOptions) (*str.WatchingResult, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/watching", *user)
	url, err := uri.AddQuery(url, options)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch watching url:" + url)
	req, err := u.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.WatchingResult)
	resp, err := u.client.Do(ctx, req, &result)

	if err != nil {
		printer.Println("fetch watching err:" + err.Error())
		return nil, resp, err
	}

	return result, resp, nil
}

// Report Report a user for moderator review.
// Send a reason and optional message with additional context.
// A user can only have one pending report per reported user.
// API docs:https://trakt.docs.apiary.io/#reference/users/report/report-a-user
func (u *UsersService) Report(ctx context.Context, user *string, report *str.UserReport) (*str.UserReportResult, *str.Response, error) {
	var url string
	url = fmt.Sprintf("users/%s/report", *user)
	printer.Println("user report")
	req, err := u.client.NewRequest(http.MethodPost, url, report)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.UserReportResult)
	resp, err := u.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusBadRequest {
		return result, resp, errors.New(*result.Message)
	}

	if resp.StatusCode == http.StatusNotFound {
		return result, resp, errors.New(*result.Message)
	}

	if err != nil {
		return nil, resp, errors.New(*result.Message)
	}

	return result, resp, nil
}
