// Package str used for structs
package str

// SavedFilterAdd represents JSON saved filter to add: a name and the Trakt URL with the filters
type SavedFilterAdd struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

func (s SavedFilterAdd) String() string {
	return Stringify(s)
}
