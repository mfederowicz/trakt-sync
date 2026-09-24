// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// SearchAddRecentHandler struct for handler
type SearchAddRecentHandler struct{}

// Handle to handle search: add_recent action
func (SearchAddRecentHandler) Handle(options *str.Options, client *internal.Client) error {
	search, err := buildRecentSearch(options)
	if err != nil {
		return err
	}

	if _, err := client.Search.AddRecentSearch(client.BuildCtxFromOptions(options), search); err != nil {
		return fmt.Errorf("add recent search error: %w", err)
	}

	printer.Printf("added %s %d for query %q to global search trends\n", search.Type, search.ID, search.Query)
	return nil
}
