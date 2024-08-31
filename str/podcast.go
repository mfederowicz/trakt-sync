// Package str used for structs
package str

// Podcast represents JSON podcast object
type Podcast struct {
	Title *string `json:"title,omitempty"`
	Year  *int    `json:"year,omitempty"`
	IDs   *IDs    `json:"ids,omitempty"`
}

func (p Podcast) String() string {
	return Stringify(p)
}
