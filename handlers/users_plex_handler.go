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

// UsersPlexSettingsHandler struct for handler
type UsersPlexSettingsHandler struct{}

// Handle to handle users: plex_settings action
func (UsersPlexSettingsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns your Plex connection, webhook, sync selection and toggles.")
	result, resp, err := client.Users.GetPlexSettings(client.BuildCtxFromOptions(options))
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}

	return writeResult(options, result)
}

// UsersUpdatePlexSettingsHandler struct for handler
type UsersUpdatePlexSettingsHandler struct{ common CommonLogic }

// Handle to handle users: update_plex_settings action
func (h UsersUpdatePlexSettingsHandler) Handle(options *str.Options, client *internal.Client) error {
	settings := new(str.PlexSettingsUpdate)
	if err := readStrictInput(&h.common, options, settings); err != nil {
		return err
	}
	if settings.Sync == nil && settings.Scrobbler == nil && settings.Webhook == nil && settings.TriggerSync == nil {
		return errors.New("plex settings update needs at least one of sync, scrobbler, webhook, trigger_sync")
	}

	printer.Println("Update your Plex settings.")
	resp, err := client.Users.UpdatePlexSettings(client.BuildCtxFromOptions(options), settings)
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}

	printer.Println("result: success, Plex settings updated")
	return nil
}

// UsersPlexConnectHandler struct for handler
type UsersPlexConnectHandler struct{}

// Handle to handle users: plex_connect action
func (UsersPlexConnectHandler) Handle(options *str.Options, client *internal.Client) error {
	returnURL := options.ReturnURL
	if len(returnURL) == consts.ZeroValue {
		returnURL = consts.DefaultReturnURL
	}

	printer.Println("Create a Plex web auth URL.")
	result, resp, err := client.Users.ConnectPlex(client.BuildCtxFromOptions(options), &str.PlexConnect{ReturnURL: &returnURL})
	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("plex_connect rejected return_url %s (allowed: trakt://..., http(s)://localhost, https://*.trakt.tv)", returnURL)
	}
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}
	if result.URL == nil || len(*result.URL) == consts.ZeroValue {
		return errors.New("plex_connect: no web auth URL in the response")
	}

	printer.Println("Open this URL to connect Plex: " + *result.URL)
	return writeResult(options, result)
}

// UsersPlexDisconnectHandler struct for handler
type UsersPlexDisconnectHandler struct{}

// Handle to handle users: plex_disconnect action
func (UsersPlexDisconnectHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Disconnect Plex: clears the authorization, selection and sync state.")
	resp, err := client.Users.DisconnectPlex(client.BuildCtxFromOptions(options))
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}

	printer.Println("result: success, Plex disconnected")
	return nil
}

// UsersPlexServersHandler struct for handler
type UsersPlexServersHandler struct{}

// Handle to handle users: plex_servers action
func (UsersPlexServersHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns your Plex servers.")
	result, resp, err := client.Users.GetPlexServers(client.BuildCtxFromOptions(options))
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}
	if len(result.Servers) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result.Servers))
	return writeResult(options, result)
}

// UsersPlexServerHandler struct for handler
type UsersPlexServerHandler struct{}

// Handle to handle users: plex_server action
func (UsersPlexServerHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New("set Plex server id ie: -i <id from users -a plex_servers>")
	}

	printer.Println("Returns home accounts and libraries of Plex server: " + options.ID)
	result, resp, err := client.Users.GetPlexServerAccounts(client.BuildCtxFromOptions(options), &options.ID)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found Plex server:%s", options.ID)
	}
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}

	return writeResult(options, result)
}

// UsersPlexSyncHandler struct for handler
type UsersPlexSyncHandler struct{}

// Handle to handle users: plex_sync action
func (UsersPlexSyncHandler) Handle(options *str.Options, client *internal.Client) error {
	request := &str.PlexSyncRequest{}
	target := "every selected server"
	if len(options.ID) > consts.ZeroValue {
		request.ServerID = &options.ID
		target = "server " + options.ID
	}
	if options.AllData {
		request.AllData = &options.AllData
	}

	printer.Println("Queue a Plex sync of " + target)
	resp, err := client.Users.SyncPlex(client.BuildCtxFromOptions(options), request)
	if resp != nil && resp.StatusCode == http.StatusUnprocessableEntity {
		return errors.New("plex_sync: no Plex server to sync (select one with users -a update_plex_settings)")
	}
	if err = plexError(options.Action, resp, err); err != nil {
		return err
	}

	printer.Println("result: success, Plex sync queued for " + target)
	return nil
}

// plexError maps a Plex settings response to a readable error. A 401 with an error_code is Plex's own
// auth failure (e.g. bad_auth); a 401 without it means Trakt does not open the route to API apps.
// Client.Do returns no error for statuses it has no type for (502, 503, 504), so the status is checked too.
func plexError(action string, resp *str.Response, err error) error {
	var invalidUser *internal.InvalidUserError
	if errors.As(err, &invalidUser) && len(invalidUser.ErrorCode) > consts.ZeroValue {
		return fmt.Errorf("%s: Plex %s: %s %s", action, invalidUser.ErrorCode, invalidUser.Message, invalidUser.Guidance)
	}
	if apiErr := notOpenToAPIApps(action, err); apiErr != nil {
		return apiErr
	}
	if err != nil {
		return fmt.Errorf("%s error: %w", action, err)
	}
	if resp != nil && resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%s: Plex request failed with status %d", action, resp.StatusCode)
	}
	return nil
}
