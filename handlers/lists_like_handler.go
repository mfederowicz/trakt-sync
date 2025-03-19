// Package handlers used to handle module actions
package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// ListsLikeHandler struct for handler
type ListsLikeHandler struct{}

// Handle to handle lists: like action
func (h ListsLikeHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyListIDMsg)
	}

	resp, _ := h.likeSingleList(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found list for:%d", options.TraktID)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Print("result: success \n")
	}

	return nil
}

func (ListsLikeHandler) likeSingleList(client *internal.Client, options *str.Options) (*str.Response, error) {
	listID := options.InternalID
	
	if !options.Remove {
		resp, err := client.Lists.LikeList(
			context.Background(),
			&listID,
		)
		return resp, err
	}
	
	resp, err := client.Lists.RemoveLikeList(
		context.Background(),
		&listID,
	)

	return resp, err
}
