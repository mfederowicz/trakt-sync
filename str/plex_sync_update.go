// Package str used for structs
package str

// PlexSyncUpdate represents JSON Plex sync selection and toggles to change
type PlexSyncUpdate struct {
	Selection *PlexSelection    `json:"selection,omitempty"`
	Toggles   *PlexToggleGroups `json:"toggles,omitempty"`
}

func (p PlexSyncUpdate) String() string {
	return Stringify(p)
}
