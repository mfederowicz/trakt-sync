// Package handlers used to handle module actions
package handlers

import (
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsPopularHandler struct for handler
type ListsPopularHandler struct{}

// Handle to handle lists: popular action
func (h ListsPopularHandler) Handle(options *str.Options, client *internal.Client) error {
	return nil
}
