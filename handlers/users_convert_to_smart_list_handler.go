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

// UsersConvertToSmartListHandler struct for handler
type UsersConvertToSmartListHandler struct{ common CommonLogic }

// Handle to handle users: convert_to_smart_list action
func (h UsersConvertToSmartListHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.SmartListDraftListIDMsg)
	}
	list, err := readSmartListWrite(&h.common, options)
	if err != nil {
		return err
	}
	if err := validSmartListWrite(list, true); err != nil {
		return err
	}

	printer.Println("Turn list " + options.ID + " into smart list: " + *list.Name)
	result, resp, err := client.Users.ConvertListToSmartList(cli.ContextFromOptions(options), options.UserName, options.ID, list)
	// smart lists are VIP Enhanced: a non-VIP account over its limit gets 420
	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf(consts.ListNotConvertibleMsg, options.ID)
	}
	if err != nil {
		return fmt.Errorf("fetch %s error: %w", options.Action, err)
	}

	if resp != nil && resp.StatusCode == http.StatusOK {
		printer.Println("result: the list was already converted")
	}
	if result.IDs != nil && result.IDs.Slug != nil {
		printer.Println("result: success, smart list:" + *result.IDs.Slug)
	}
	return writeSmartList(options, result)
}
