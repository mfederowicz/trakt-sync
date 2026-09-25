// Package str used for structs
package str

// ReviewThanks represents JSON popular shows and movies the user hasn't watched yet
type ReviewThanks struct {
	Shows  []*MediaItem `json:"shows,omitempty"`
	Movies []*MediaItem `json:"movies,omitempty"`
}

func (r ReviewThanks) String() string {
	return Stringify(r)
}
