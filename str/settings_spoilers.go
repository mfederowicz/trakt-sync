// Package str used for structs
package str

// SettingsSpoilers represents JSON spoiler settings; each value is show or hide
type SettingsSpoilers struct {
	Episodes *string `json:"episodes,omitempty"`
	Shows    *string `json:"shows,omitempty"`
	Movies   *string `json:"movies,omitempty"`
}

func (s SettingsSpoilers) String() string {
	return Stringify(s)
}
