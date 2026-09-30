// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// SearchRemoveRecentHandler struct for handler
type SearchRemoveRecentHandler struct{}

// Handle to handle search: remove_recent action
func (SearchRemoveRecentHandler) Handle(options *str.Options, client *trakt.Client) error {
	search, err := buildRecentSearch(options)
	if err != nil {
		return err
	}

	if _, err := client.Search.RemoveRecentSearch(client.BuildCtxFromOptions(options), search); err != nil {
		return fmt.Errorf("remove recent search error: %w", err)
	}

	printer.Printf("removed %s %d for query %q from global search trends\n", search.Type, search.ID, search.Query)
	return nil
}
