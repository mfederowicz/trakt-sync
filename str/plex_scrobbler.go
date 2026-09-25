// Package str used for structs
package str

// PlexScrobbler represents JSON Plex real-time scrobbler toggles
type PlexScrobbler struct {
	Toggles *PlexToggleGroups `json:"toggles,omitempty"`
}

func (p PlexScrobbler) String() string {
	return Stringify(p)
}
