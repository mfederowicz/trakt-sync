// Package str used for structs
package str

// PlexServers represents JSON list of the user's Plex servers
type PlexServers struct {
	Servers []*PlexServer `json:"servers,omitempty"`
}

func (p PlexServers) String() string {
	return Stringify(p)
}
