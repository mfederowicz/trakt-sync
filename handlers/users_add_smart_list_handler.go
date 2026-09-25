// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersAddSmartListHandler struct for handler
type UsersAddSmartListHandler struct{ common CommonLogic }

// Handle to handle users: add_smart_list action
func (h UsersAddSmartListHandler) Handle(options *str.Options, client *internal.Client) error {
	list, err := readSmartListWrite(&h.common, options)
	if err != nil {
		return err
	}
	if err := validSmartListWrite(list, true); err != nil {
		return err
	}

	printer.Println("Create smart list: " + *list.Name)
	result, resp, err := client.Users.AddSmartList(client.BuildCtxFromOptions(options), &options.UserName, list)
	// smart lists are VIP Enhanced: a non-VIP account over its limit gets 420
	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}
	if err = smartListError(options.Action, *list.Name, resp, err); err != nil {
		return err
	}

	if result.IDs != nil && result.IDs.Slug != nil {
		printer.Println("result: success, created smart list:" + *result.IDs.Slug)
	}
	return writeSmartList(options, result)
}
