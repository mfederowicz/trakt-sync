// Package str used for structs
package str

// PlexConnectResult represents JSON Plex connect response with the web auth URL
type PlexConnectResult struct {
	URL *string `json:"url,omitempty"`
}

func (p PlexConnectResult) String() string {
	return Stringify(p)
}
