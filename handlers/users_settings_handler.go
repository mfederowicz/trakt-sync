// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersSettingsHandler struct for handler
type UsersSettingsHandler struct{}

// Handle to handle users: settings action
func (UsersSettingsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("users settings handler:" + options.UserName)

	settings, _, err := fetchUsersSettings(client, options)
	if err != nil {
		return fmt.Errorf("fetch settings error:%w", err)
	}

	printer.Print("Found " + options.Action + " data \n")
	jsonData, err := json.MarshalIndent(settings, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)
	return nil
}

func fetchUsersSettings(client *trakt.Client, options *str.Options) (*str.UserSettings, *str.Response, error) {
	settings, resp, err := client.Users.GetSettings(cli.ContextFromOptions(options))

	return settings, resp, err
}
