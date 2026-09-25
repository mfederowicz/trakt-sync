// Package str used for structs
package str

// PlexServer represents JSON Plex server; url is null when the server is unreachable
type PlexServer struct {
	ID                *string `json:"id,omitempty"`
	Name              *string `json:"name,omitempty"`
	ConnectionCount   *int    `json:"connection_count,omitempty"`
	ConnectionTimeout *int    `json:"connection_timeout,omitempty"`
	Ports             []int   `json:"ports,omitempty"`
	Owned             *bool   `json:"owned,omitempty"`
	URL               *string `json:"url,omitempty"`
}

func (p PlexServer) String() string {
	return Stringify(p)
}
