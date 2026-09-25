// Package str used for structs
package str

// PlexLibrarySelection represents JSON selected Plex library
type PlexLibrarySelection struct {
	ServerID *string `json:"server_id,omitempty"`
	UUID     *string `json:"uuid,omitempty"`
}

func (p PlexLibrarySelection) String() string {
	return Stringify(p)
}
