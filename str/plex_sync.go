// Package str used for structs
package str

// PlexSync represents JSON Plex batch sync state, selection and toggles
type PlexSync struct {
	Configured  *bool             `json:"configured,omitempty"`
	Error       *bool             `json:"error,omitempty"`
	ServerLimit *int              `json:"server_limit,omitempty"`
	Selection   *PlexSelection    `json:"selection,omitempty"`
	Toggles     *PlexToggleGroups `json:"toggles,omitempty"`
}

func (p PlexSync) String() string {
	return Stringify(p)
}
