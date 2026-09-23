// Package str used for structs
package str

// MediaItem represents JSON media item object: a movie or a show with its stats
type MediaItem struct {
	Watchers  *int   `json:"watchers,omitempty"`
	ListCount *int   `json:"list_count,omitempty"`
	Movie     *Movie `json:"movie,omitempty"`
	Show      *Show  `json:"show,omitempty"`
}

func (m MediaItem) String() string {
	return Stringify(m)
}
