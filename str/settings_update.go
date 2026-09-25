// Package str used for structs
package str

// SettingsUpdate represents JSON users/settings update request object; send only the values to change
type SettingsUpdate struct {
	User     *SettingsUpdateUser `json:"user,omitempty"`
	Browsing *SettingsBrowsing   `json:"browsing,omitempty"`
}

func (s SettingsUpdate) String() string {
	return Stringify(s)
}
