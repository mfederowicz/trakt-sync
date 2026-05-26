// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersUpdateListItemHandler struct for handler
type UsersUpdateListItemHandler struct{ common CommonLogic }

// Handle to handle sync: update_list_item action
func (m UsersUpdateListItemHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}
	if options.ListItemID == consts.ZeroValue {
		return errors.New(consts.EmptyListItemIDMsg)
	}
	resp, err := m.usersUpdateListItem(client, options)
	if err != nil {
		return fmt.Errorf("update personal list item error:%w", err)
	}
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("update personal list item success for list item:", options.ListItemID)
	}

	return nil
}

func (UsersUpdateListItemHandler) usersUpdateListItem(client *internal.Client, options *str.Options) (*str.Response, error) {
	item := new(str.PersonalListItem)
	if len(options.Notes) > consts.ZeroValue {
		item.Notes = &options.Notes
	}

	resp, err := client.Users.UpdateListItem(
		client.BuildCtxFromOptions(options),
		&options.UserName,
		&options.ID,
		&options.ListItemID,
		item)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("list item not found:%d", options.ListItemID)
	}

	return resp, err
}
