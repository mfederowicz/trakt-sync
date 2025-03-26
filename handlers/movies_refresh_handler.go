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

// MoviesRefreshHandler struct for handler
type MoviesRefreshHandler struct{}

// Handle to handle movie: refresh action
func (h MoviesRefreshHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Queue this movie for a full metadata and image refresh.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	resp, _ := h.refreshMovie(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found movie for:%s", options.InternalID)
	}

	if resp.StatusCode == http.StatusUpgradeRequired {
		return cli.HandleUpgrade(resp)
	}

	if resp.StatusCode == http.StatusConflict {
		return errors.New("result: movie is already queued")
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Print("result: success \n")
	}

	return nil
}

func (MoviesRefreshHandler) refreshMovie(client *internal.Client, options *str.Options) (*str.Response, error) {
	movieID := options.InternalID
	resp, err := client.Movies.RefreshMovieMetadata(
		context.Background(),
		&movieID,
	)
	return resp, err
}
