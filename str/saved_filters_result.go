// Package str used for structs
package str

// SavedFiltersResult represents JSON add saved filters response
type SavedFiltersResult struct {
	Added   []*SavedFilter    `json:"added,omitempty"`
	Skipped []*SavedFilterAdd `json:"skipped,omitempty"`
}

func (s SavedFiltersResult) String() string {
	return Stringify(s)
}
