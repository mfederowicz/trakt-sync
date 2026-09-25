// Package str used for structs
package str

// PlexSelection represents JSON selected Plex servers, libraries and home users
type PlexSelection struct {
	ServerIDs  []string                `json:"server_ids,omitempty"`
	LibraryIDs []*PlexLibrarySelection `json:"library_ids,omitempty"`
	UserIDs    []string                `json:"user_ids,omitempty"`
}

func (p PlexSelection) String() string {
	return Stringify(p)
}
