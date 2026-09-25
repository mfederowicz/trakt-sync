// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersDeleteSavedFilterHandler struct for handler
type UsersDeleteSavedFilterHandler struct{}

// Handle to handle users: delete_saved_filter action
func (UsersDeleteSavedFilterHandler) Handle(options *str.Options, client *internal.Client) error {
	id, err := strconv.Atoi(options.ID)
	if err != nil || id <= consts.ZeroValue {
		return errors.New("set saved filter id ie: -i 101 (ids from users -a saved_filters)")
	}

	printer.Printf("Delete saved filter %d (VIP only)\n", id)
	resp, err := client.Users.DeleteSavedFilter(client.BuildCtxFromOptions(options), id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found saved filter:%d", id)
	}
	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}
	if err != nil {
		return fmt.Errorf("delete saved filter error: %w", err)
	}

	printer.Printf("result: success, deleted saved filter:%d\n", id)
	return nil
}
