// Package handlers used to handle module actions
package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// PeopleRefreshHandler struct for handler
type PeopleRefreshHandler struct{}

// Handle to handle people: refresh action
func (h PeopleRefreshHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyPersonIDMsg)
	}

	resp, _ := h.refreshPerson(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found person for:%s", options.ID)
	}

	if resp.StatusCode == http.StatusUpgradeRequired {
		return cli.HandleUpgrade(resp)
	}

	if resp.StatusCode == http.StatusConflict {
		return errors.New("result: person is already queued")
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Print("result: success \n")
	}

	return nil
}

func (PeopleRefreshHandler) refreshPerson(client *internal.Client, options *str.Options) (*str.Response, error) {
	personID := options.ID
	resp, err := client.People.RefreshPersonMetadata(
		context.Background(),
		&personID,
	)
	return resp, err
}
