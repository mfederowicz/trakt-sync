// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersUpdateSettingsHandler struct for handler
type UsersUpdateSettingsHandler struct{ common CommonLogic }

// Handle to handle users: update_settings action
func (h UsersUpdateSettingsHandler) Handle(options *str.Options, client *internal.Client) error {
	settings := new(str.SettingsUpdate)
	if err := readStrictInput(&h.common, options, settings); err != nil {
		return err
	}
	if err := validSettingsUpdate(settings); err != nil {
		return err
	}

	printer.Println("Update settings of the authenticated user.")
	if _, err := client.Users.UpdateSettings(client.BuildCtxFromOptions(options), settings); err != nil {
		return fmt.Errorf("update settings error: %w", err)
	}

	printer.Println("result: success, settings updated")
	return nil
}

// validSettingsUpdate checks that the update changes something and that the enum values are contract values.
func validSettingsUpdate(settings *str.SettingsUpdate) error {
	if settings.User == nil && settings.Browsing == nil {
		return errors.New("settings update needs a user or browsing block")
	}
	if settings.Browsing == nil {
		return nil
	}

	checks := []enumCheck{
		{name: "browsing.dark_knight", value: settings.Browsing.DarkKnight, valid: cfg.SettingsDarkKnightValues},
	}
	if spoilers := settings.Browsing.Spoilers; spoilers != nil {
		checks = append(checks,
			enumCheck{name: "browsing.spoilers.episodes", value: spoilers.Episodes, valid: cfg.SettingsSpoilerValues},
			enumCheck{name: "browsing.spoilers.shows", value: spoilers.Shows, valid: cfg.SettingsSpoilerValues},
			enumCheck{name: "browsing.spoilers.movies", value: spoilers.Movies, valid: cfg.SettingsSpoilerValues},
		)
	}

	return checkEnums(checks)
}
