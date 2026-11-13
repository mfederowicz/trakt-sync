// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// SmartListsTemplatesHandler struct for handler
type SmartListsTemplatesHandler struct{}

// Handle to handle smart_lists: templates action
func (SmartListsTemplatesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns the official smart list templates")
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.SmartLists.GetSmartListTemplates(cli.ContextFromOptions(options), &opts)
	if err != nil {
		return fmt.Errorf("fetch smart list templates error: %w", err)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeSmartList(options, result)
}
