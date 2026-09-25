// Package str used for structs
package str

// SettingsBrowsing represents JSON browsing values of a settings update
type SettingsBrowsing struct {
	Genres           *SettingsGenres   `json:"genres,omitempty"`
	Spoilers         *SettingsSpoilers `json:"spoilers,omitempty"`
	WatchNow         *SettingsWatchNow `json:"watchnow,omitempty"`
	DarkKnight       *string           `json:"dark_knight,omitempty"`
	WatchOnlyOnce    *bool             `json:"watch_only_once,omitempty"`
	ShowRatingPrompt *bool             `json:"show_rating_prompt,omitempty"`
	Locale           *string           `json:"locale,omitempty"`
}

func (s SettingsBrowsing) String() string {
	return Stringify(s)
}
