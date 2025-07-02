// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// SyncRemovePlaybackHandler struct for handler
type SyncRemovePlaybackHandler struct{ common CommonLogic }

// Handle to handle sync: remove_playback action
func (m SyncRemovePlaybackHandler) Handle(options *str.Options, client *internal.Client) error {
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

func (SyncRemovePlaybackHandler) syncRemovePlaybackItem(client *internal.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Sync.RemovePlaybackItem(
		client.BuildCtxFromOptions(options),
		&options.PlaybackID,
	)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("result: success, remove playback item:%d", options.PlaybackID)
	}

	return nil, nil
}
