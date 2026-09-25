// Package str used for structs
package str

// SettingsWatchNow represents JSON watch now settings
type SettingsWatchNow struct {
	Country       *string  `json:"country,omitempty"`
	Favorites     []string `json:"favorites,omitempty"`
	OnlyFavorites *bool    `json:"only_favorites,omitempty"`
}

func (s SettingsWatchNow) String() string {
	return Stringify(s)
}
