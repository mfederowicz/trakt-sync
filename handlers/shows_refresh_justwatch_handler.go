// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// ShowsRefreshJustwatchHandler struct for handler
type ShowsRefreshJustwatchHandler struct{}

// Handle to handle shows: refresh_justwatch action
func (ShowsRefreshJustwatchHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	printer.Println("Queue this show for a JustWatch links refresh (VIP only).")
	resp, err := client.Shows.RefreshShowJustwatch(client.BuildCtxFromOptions(options), &options.InternalID)
	if resp == nil {
		return fmt.Errorf("refresh justwatch error: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found show for:%s", options.InternalID)
	}

	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}

	if err != nil {
		return fmt.Errorf("refresh justwatch error: %w", err)
	}

	printer.Println("result: success")
	return nil
}
