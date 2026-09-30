// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// PeopleRefreshHandler struct for handler
type PeopleRefreshHandler struct{}

// Handle to handle people: refresh action
func (h PeopleRefreshHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyPersonIDMsg)
	}

	resp, err := h.refreshPerson(client, options)
	if resp == nil {
		return fmt.Errorf("refresh person error: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found person for:%s", options.ID)
	}

	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}

	if resp.StatusCode == http.StatusConflict {
		return errors.New("result: person is already queued")
	}

	if err != nil {
		return fmt.Errorf("refresh person error: %w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Print("result: success \n")
	}

	return nil
}

func (PeopleRefreshHandler) refreshPerson(client *trakt.Client, options *str.Options) (*str.Response, error) {
	personID := options.ID
	resp, err := client.People.RefreshPersonMetadata(
		cli.ContextFromOptions(options),
		personID,
	)
	return resp, err
}
