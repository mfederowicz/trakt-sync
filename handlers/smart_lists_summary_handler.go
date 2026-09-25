// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// SmartListsSummaryHandler struct for handler
type SmartListsSummaryHandler struct{}

// Handle to handle smart_lists: summary action
func (SmartListsSummaryHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validSmartListID(options); err != nil {
		return err
	}

	printer.Println("Returns the smart list definition: " + options.InternalID)
	result, resp, err := client.SmartLists.GetSmartList(client.BuildCtxFromOptions(options), &options.InternalID)
	if err = smartListError(options.Action, options.InternalID, resp, err); err != nil {
		return err
	}
	if isEmptySmartList(result) {
		return fmt.Errorf(consts.SmartListNotFoundMsg, options.InternalID)
	}

	return writeSmartList(options, result)
}
