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

// UsersSmartListDraftHandler struct for handler
type UsersSmartListDraftHandler struct{}

// Handle to handle users: smart_list_draft action
func (UsersSmartListDraftHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.SmartListDraftListIDMsg)
	}

	printer.Println("Returns the smart list draft of list: " + options.ID)
	result, resp, err := client.Users.GetSmartListDraft(cli.ContextFromOptions(options), options.UserName, options.ID)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf(consts.SmartListDraftNotFoundMsg, options.ID)
	}
	if err != nil {
		return fmt.Errorf("fetch %s error: %w", options.Action, err)
	}
	if result == nil || (result.List == nil && result.Source == nil) {
		return fmt.Errorf(consts.SmartListDraftNotFoundMsg, options.ID)
	}

	return writeResult(options, result)
}
