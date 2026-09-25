// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersUpdateSmartListHandler struct for handler
type UsersUpdateSmartListHandler struct{ common CommonLogic }

// Handle to handle users: update_smart_list action
func (h UsersUpdateSmartListHandler) Handle(options *str.Options, client *internal.Client) error {
	options.InternalID = options.ID
	if err := validSmartListID(options); err != nil {
		return err
	}
	list, err := readSmartListWrite(&h.common, options)
	if err != nil {
		return err
	}
	if err := validSmartListWrite(list, false); err != nil {
		return err
	}

	printer.Println("Update smart list: " + options.ID)
	result, resp, err := client.Users.UpdateSmartList(client.BuildCtxFromOptions(options), &options.UserName, &options.ID, list)
	if err = smartListError(options.Action, options.ID, resp, err); err != nil {
		return err
	}

	return writeSmartList(options, result)
}
