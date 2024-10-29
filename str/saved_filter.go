// Package str used for structs
package str

// SavedFilter represents JSON filter object
type SavedFilter struct {
	Rank      *int       `json:"rank,omitempty"`
	ID        *int64     `json:"id,omitempty"`
	Section   *string    `json:"section,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Path      *string    `json:"path,omitempty"`
	Query     *string    `json:"query,omitempty"`
	UpdatedAt *Timestamp `json:"updated_at,omitempty"`
}

func (i SavedFilter) String() string {
	return Stringify(i)
}
