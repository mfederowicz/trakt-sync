// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersSmartListHandler struct for handler
type UsersSmartListHandler struct{}

// Handle to handle users: smart_list action
func (UsersSmartListHandler) Handle(options *str.Options, client *internal.Client) error {
	options.InternalID = options.ID
	if err := validSmartListID(options); err != nil {
		return err
	}

	printer.Println("Returns the smart list definition: " + options.ID)
	result, resp, err := client.Users.GetSmartList(client.BuildCtxFromOptions(options), &options.UserName, &options.ID)
	if err = smartListError(options.Action, options.ID, resp, err); err != nil {
		return err
	}
	if isEmptySmartList(result) {
		return fmt.Errorf(consts.SmartListNotFoundMsg, options.ID)
	}

	return writeSmartList(options, result)
}
