// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncPlaybackHandler struct for handler
type SyncPlaybackHandler struct{ common CommonLogic }

// Handle to handle sync: playback action
func (m SyncPlaybackHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns playback progress.")

	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	err = m.common.CheckDates(options.StartDate, options.EndDate, options.Timezone)
	if err != nil {
		return err
	}

	result, _, err := m.syncPlayback(client, options, consts.DefaultPage)

	if err != nil {
		return err
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (m SyncPlaybackHandler) syncPlayback(client *internal.Client, options *str.Options, page int) ([]*str.PlaybackProgress, *str.Response, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, StartAt: options.StartDate, EndAt: options.EndDate}
	list, resp, err := client.Sync.GetPlaybackProgress(
		client.BuildCtxFromOptions(options),
		&options.Type,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, _, err := m.syncPlayback(client, options, nextPage)
		if err != nil {
			return nil, nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, resp, nil
}
