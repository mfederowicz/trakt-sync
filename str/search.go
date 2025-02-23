// Package str used for structs
package str

// Search represents JSON search object
type Search struct {
	RecentCount *int `json:"recent_count,omitempty"`
}

func (s Search) String() string {
	return Stringify(s)
}
