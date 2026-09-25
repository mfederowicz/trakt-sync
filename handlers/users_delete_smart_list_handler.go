// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersDeleteSmartListHandler struct for handler
type UsersDeleteSmartListHandler struct{}

// Handle to handle users: delete_smart_list action
func (UsersDeleteSmartListHandler) Handle(options *str.Options, client *internal.Client) error {
	options.InternalID = options.ID
	if err := validSmartListID(options); err != nil {
		return err
	}

	printer.Println("Delete smart list: " + options.ID)
	resp, err := client.Users.DeleteSmartList(client.BuildCtxFromOptions(options), &options.UserName, &options.ID)
	if err = smartListError(options.Action, options.ID, resp, err); err != nil {
		return err
	}

	printer.Println("result: success, deleted smart list:" + options.ID)
	return nil
}
