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

// ShowsRefreshHandler struct for handler
type ShowsRefreshHandler struct{}

// Handle to handle show: refresh action
func (h ShowsRefreshHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Queue this show for a full metadata and image refresh.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	resp, _ := h.refreshShow(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found show for:%s", options.InternalID)
	}

	if resp.StatusCode == http.StatusUpgradeRequired {
		return cli.HandleUpgrade(resp)
	}

	if resp.StatusCode == http.StatusConflict {
		return errors.New("result: show is already queued")
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Print("result: success \n")
	}

	return nil
}

func (ShowsRefreshHandler) refreshShow(client *internal.Client, options *str.Options) (*str.Response, error) {
	showID := options.InternalID
	resp, err := client.Shows.RefreshShowMetadata(
		client.BuildCtxFromOptions(options),
		&showID,
	)
	return resp, err
}
