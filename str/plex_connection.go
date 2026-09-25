// Package str used for structs
package str

// PlexConnection represents JSON Plex connection status
type PlexConnection struct {
	Connected *bool   `json:"connected,omitempty"`
	Username  *string `json:"username,omitempty"`
}

func (p PlexConnection) String() string {
	return Stringify(p)
}
