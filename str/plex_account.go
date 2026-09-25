// Package str used for structs
package str

// PlexAccount represents JSON Plex home account
type PlexAccount struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func (p PlexAccount) String() string {
	return Stringify(p)
}
