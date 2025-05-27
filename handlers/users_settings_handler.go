// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersSettingsHandler struct for handler
type UsersSettingsHandler struct{}

// Handle to handle users: settings action
func (UsersSettingsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("users settings handler:" + options.UserName)

	settings, _, err := fetchUsersSettings(client, options)
	if err != nil {
		return fmt.Errorf("fetch settings error:%w", err)
	}

	printer.Print("Found " + options.Action + " data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(settings, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchUsersSettings(client *internal.Client, options *str.Options) (*str.UserSettings, *str.Response, error) {
	settings, resp, err := client.Users.RetrieveSettings(client.BuildCtxFromOptions(options))

	return settings, resp, err
}
