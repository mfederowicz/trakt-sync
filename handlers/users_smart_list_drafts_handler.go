// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// UsersSmartListDraftsHandler struct for handler
type UsersSmartListDraftsHandler struct{}

// Handle to handle users: smart_list_drafts action
func (UsersSmartListDraftsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns the lists ready to become smart lists of: " + options.UserName)
	result, _, err := client.Users.GetSmartListDrafts(cli.ContextFromOptions(options), options.UserName)
	if err != nil {
		return fmt.Errorf("fetch %s error: %w", options.Action, err)
	}

	count := consts.ZeroValue
	if result.ConvertibleCount != nil {
		count = *result.ConvertibleCount
	}
	printer.Printf("Found %d result \n", count)
	return writeResult(options, result)
}
