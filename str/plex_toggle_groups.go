// Package str used for structs
package str

// PlexToggleGroups represents JSON Plex toggles per media type
type PlexToggleGroups struct {
	Movie   *PlexToggles `json:"movie,omitempty"`
	Show    *PlexToggles `json:"show,omitempty"`
	Season  *PlexToggles `json:"season,omitempty"`
	Episode *PlexToggles `json:"episode,omitempty"`
}

func (p PlexToggleGroups) String() string {
	return Stringify(p)
}
