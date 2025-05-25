// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersStatsHandler struct for handler
type UsersStatsHandler struct{}

// Handle to handle users: stats action
func (UsersStatsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("fetch stats for:" + options.UserName)
	stats, resp, err := fetchUsersStats(client, options)
	if err != nil {
		return fmt.Errorf("fetch user stats error:%w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found user stats for:%s", options.UserName)
	}

	printer.Printf("Found %s user stats\n", options.UserName)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(stats, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchUsersStats(client *internal.Client, options *str.Options) (*str.UserStats, *str.Response, error) {
	username := options.UserName
	stats, resp, err := client.Users.GetStats(
		client.BuildCtxFromOptions(options),
		&username,
	)

	return stats, resp, err
}
