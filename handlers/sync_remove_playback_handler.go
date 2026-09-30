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

// SyncRemovePlaybackHandler struct for handler
type SyncRemovePlaybackHandler struct{ common CommonLogic }

// Handle to handle sync: remove_playback action
func (m SyncRemovePlaybackHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.PlaybackID == consts.ZeroValue {
		return errors.New("empty playback_id")
	}

	printer.Println("Remove playback item:", options.PlaybackID)
	_, err := m.syncRemovePlaybackItem(client, options)

	if err != nil {
		return err
	}

	return nil
}

func (SyncRemovePlaybackHandler) syncRemovePlaybackItem(client *trakt.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Sync.RemovePlaybackItem(
		cli.ContextFromOptions(options),
		options.PlaybackID,
	)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, remove playback item:%d\n", options.PlaybackID)
	}

	return nil, nil
}
