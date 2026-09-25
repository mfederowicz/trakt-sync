// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersAddSavedFiltersHandler struct for handler
type UsersAddSavedFiltersHandler struct{ common CommonLogic }

// Handle to handle users: add_saved_filters action
func (h UsersAddSavedFiltersHandler) Handle(options *str.Options, client *internal.Client) error {
	filters := []*str.SavedFilterAdd{}
	if err := readStrictInput(&h.common, options, &filters); err != nil {
		return err
	}
	if len(filters) == consts.ZeroValue {
		return errors.New("add_saved_filters needs a JSON array of {\"name\": ..., \"url\": ...}")
	}
	for i, f := range filters {
		if f == nil || f.Name == nil || len(*f.Name) == consts.ZeroValue || f.URL == nil || len(*f.URL) == consts.ZeroValue {
			return fmt.Errorf("saved filter %d needs name and url", i+consts.OneValue)
		}
	}

	printer.Printf("Add %d saved filters (VIP only)\n", len(filters))
	result, resp, err := client.Users.AddSavedFilters(client.BuildCtxFromOptions(options), filters)
	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}
	if err != nil {
		return fmt.Errorf("add saved filters error: %w", err)
	}

	printer.Printf("result: added %d, skipped %d\n", len(result.Added), len(result.Skipped))
	return writeResult(options, result)
}
