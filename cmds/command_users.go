// Package cmds used for commands modules
package cmds

import (
	"flag"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/str"
)

var (
	username   = "me"
	exportData []*str.PersonalList

	_usersSort                   = flag.String("s", cfg.DefaultConfig().UsersSort, consts.SortUsage)
	_usersListID                 = flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
	_usersListItemID             = UsersCmd.Flag.Int("list_item_id", cfg.DefaultConfig().ListItemID, consts.ListItemIDUsage)
	_usersItemID                 = UsersCmd.Flag.Int("item_id", cfg.DefaultConfig().ItemID, consts.ItemIDUsage)
	_usersNotes                  = UsersCmd.Flag.String("notes", cfg.DefaultConfig().Notes, consts.NotesUsage)
	_usersAction                 = UsersCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_usersType                   = UsersCmd.Flag.String("t", cfg.DefaultConfig().UsersType, consts.UsersTypeUsage)
	_usersSection                = UsersCmd.Flag.String("section", cfg.DefaultConfig().UsersSection, consts.UsersSectionUsage)
	_usersReason                 = UsersCmd.Flag.String("r", cfg.DefaultConfig().Reason, consts.ReasonUsage)
	_usersMessage                = UsersCmd.Flag.String("message", cfg.DefaultConfig().Msg, consts.ReportMsgUsage)
	_usersDeny                   = UsersCmd.Flag.Bool("deny", cfg.DefaultConfig().Deny, consts.DenyUsage)
	_usersDelete                 = UsersCmd.Flag.Bool("delete", cfg.DefaultConfig().Delete, consts.DeleteUsage)
	_usersPrivacy                = UsersCmd.Flag.String("privacy", cfg.DefaultConfig().Privacy, consts.PrivacyUsage)
	_usersDisplayNumbers         = UsersCmd.Flag.Bool("display_numbers", cfg.DefaultConfig().DisplayNumbers, consts.DisplayNumbersUsage)
	_usersAllowComments          = UsersCmd.Flag.Bool("allow_comments", cfg.DefaultConfig().AllowComments, consts.AllowCommentsUsage)
	_usersFollowerRequest        = UsersCmd.Flag.Int("follower_request", cfg.DefaultConfig().FollowerRequest, consts.FollowerRequestUsage)
	_usersItems                  = UsersCmd.Flag.String("items", consts.EmptyString, consts.ItemsUsage)
	_usersCommentsIncludeReplies = UsersCmd.Flag.String("include_replies", cfg.DefaultConfig().IncludeReplies, consts.IncludeRepliesUsage)
	_usersCommentsCommentType    = UsersCmd.Flag.String("comment_type", cfg.DefaultConfig().CommentType, consts.CommentTypeUsage)
)

// UsersCmd Returns all personal lists for a user.
var UsersCmd = &Command{
	Name:    "users",
	Usage:   "",
	Summary: "Returns all data for a users.",
	Help:    `users command`,
}

func usersListsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	err := cmd.ValidModuleActionType(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}
	err = cmd.ValidSection(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}
	err = cmd.ValidSort(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	var handler handlers.UsersHandler
	allHandlers := map[string]handlers.Handler{
		"settings":            handlers.UsersSettingsHandler{},
		"following_requests":  handlers.UsersFollowingRequestsHandler{},
		"follower_requests":   handlers.UsersFollowerRequestsHandler{},
		"saved_filters":       handlers.UsersSavedFiltersHandler{},
		"hidden_items":        handlers.UsersHiddenItemsHandler{},
		"add_hidden_items":    handlers.UsersAddHiddenItemsHandler{},
		"remove_hidden_items": handlers.UsersRemoveHiddenItemsHandler{},
		"profile":             handlers.UsersProfileHandler{},
		"likes":               handlers.UsersLikesHandler{},
		"collection":          handlers.UsersCollectionHandler{},
		"comments":            handlers.UsersCommentsHandler{},
		"notes":               handlers.UsersNotesHandler{},
		"lists":               handlers.UsersListsHandler{},
		"add_list":            handlers.UsersAddListHandler{},
		"reorder_lists":       handlers.UsersReorderListsHandler{},
		"collaborations":      handlers.UsersCollaborationsHandler{},
		"list":                handlers.UsersListHandler{},
		"update_list":         handlers.UsersUpdateListHandler{},
		"delete_list":         handlers.UsersDeleteListHandler{},
		"list_likes":          handlers.UsersListLikesHandler{},
		"list_like":           handlers.UsersListLikeHandler{},
		"list_items":          handlers.UsersListItemsHandler{},
		"add_list_items":      handlers.UsersAddListItemsHandler{},
		"remove_list_items":   handlers.UsersRemoveListItemsHandler{},
		"reorder_list_items":  handlers.UsersReorderListItemsHandler{},
		"update_list_item":    handlers.UsersUpdateListItemHandler{},
		"list_comments":       handlers.UsersListCommentsHandler{},
		"list_report":         handlers.UsersListReportHandler{},
		"follow":              handlers.UsersFollowHandler{},
		"unfollow":            handlers.UsersUnfollowHandler{},
		"blocked_users":       handlers.UsersBlockedUsersHandler{},
		"block":               handlers.UsersBlockHandler{},
		"unblock":             handlers.UsersUnblockHandler{},
		"followers":           handlers.UsersFollowersHandler{},
		"following":           handlers.UsersFollowingHandler{},
		"friends":             handlers.UsersFriendsHandler{},
		"history":             handlers.UsersHistoryHandler{},
		"ratings":             handlers.UsersRatingsHandler{},
		"watchlist":           handlers.UsersWatchlistHandler{},
		"watchlist_comments":  handlers.UsersWatchlistCommentsHandler{},
		"favorites":           handlers.UsersFavoritesHandler{},
		"stats":               handlers.UsersStatsHandler{},
		"watched":             handlers.UsersWatchedHandler{},
	}

	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{"settings", "following_requests", "follower_requests",
		"follow_request", "saved_filters", "hidden_items", "add_hidden_items",
		"remove_hidden_items", "profile", "likes", "collection", "comments",
		"notes", "lists", "add_list", "reorder_lists", "collaborations", "list",
		"update_list", "delete_list", "list_likes", "list_like", "list_items",
		"add_list_items", "remove_list_items", "reorder_list_items", "update_list_item",
		"list_comments", "list_report", "follow", "unfollow", "blocked_users",
		"block", "unblock", "followers", "following", "friends", "history", "ratings",
		"watchlist", "watchlist_comments", "favorites", "stats", "watched"}
	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	usersListItemsDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	UsersCmd.Run = usersListsFunc
}
