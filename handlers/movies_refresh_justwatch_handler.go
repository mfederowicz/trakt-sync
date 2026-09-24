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

// MoviesRefreshJustwatchHandler struct for handler
type MoviesRefreshJustwatchHandler struct{}

// Handle to handle movies: refresh_justwatch action
func (MoviesRefreshJustwatchHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	printer.Println("Queue this movie for a JustWatch links refresh (VIP only).")
	resp, err := client.Movies.RefreshMovieJustwatch(client.BuildCtxFromOptions(options), &options.InternalID)
	if resp == nil {
		return fmt.Errorf("refresh justwatch error: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found movie for:%s", options.InternalID)
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
