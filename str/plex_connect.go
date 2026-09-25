// Package str used for structs
package str

// PlexConnect represents JSON Plex connect request
type PlexConnect struct {
	ReturnURL *string `json:"return_url,omitempty"`
}

func (p PlexConnect) String() string {
	return Stringify(p)
}
