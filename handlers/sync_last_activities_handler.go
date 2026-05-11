// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncLastActivitiesHandler struct for handler
type SyncLastActivitiesHandler struct{}

// Handle to handle sync: last_activities action
func (m SyncLastActivitiesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get last activities")
	activities, err := m.syncLastActivities(client, options)
	if err != nil {
		return fmt.Errorf("fetch last activities error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(activities, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncLastActivitiesHandler) syncLastActivities(client *internal.Client, options *str.Options) (*str.UserLastActivities, error) {
	activities, _, err := client.Sync.GetLastActivity(
		client.BuildCtxFromOptions(options),
	)

	if err != nil {
		return nil, err
	}

	return activities, nil
}
