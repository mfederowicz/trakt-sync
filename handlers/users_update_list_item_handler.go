// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// UsersUpdateListItemHandler struct for handler
type UsersUpdateListItemHandler struct{ common CommonLogic }

// Handle to handle sync: update_list_item action
func (m UsersUpdateListItemHandler) Handle(options *str.Options, client *trakt.Client) error {
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
		printer.Println("update personal list item success for list item:", options.ListItemID)
	}

	return nil
}

func (UsersUpdateListItemHandler) usersUpdateListItem(client *trakt.Client, options *str.Options) (*str.Response, error) {
	item := new(str.PersonalListItem)
	if len(options.Notes) > consts.ZeroValue {
		item.Notes = &options.Notes
	}

	resp, err := client.Users.UpdateListItem(
		cli.ContextFromOptions(options),
		options.UserName,
		options.ID,
		options.ListItemID,
		item)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("list item not found:%d", options.ListItemID)
	}

	return resp, err
}
