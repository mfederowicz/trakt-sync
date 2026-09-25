// Package str used for structs
package str

// PlexLibrary represents JSON Plex library; selected tells if it is in the sync selection
type PlexLibrary struct {
	ID       *int    `json:"id,omitempty"`
	UUID     *string `json:"uuid,omitempty"`
	Type     *string `json:"type,omitempty"`
	Title    *string `json:"title,omitempty"`
	Agent    *string `json:"agent,omitempty"`
	Scanner  *string `json:"scanner,omitempty"`
	Selected *bool   `json:"selected,omitempty"`
	URL      *string `json:"url,omitempty"`
}

func (p PlexLibrary) String() string {
	return Stringify(p)
}
