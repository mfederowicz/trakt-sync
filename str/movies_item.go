// Package str used for structs
package str

// MoviesItem represents JSON movies item object
type MoviesItem struct {
	Revenue        *int   `json:"revenue,omitempty"`
	UserCount      *int   `json:"user_count,omitempty"`
	WatcherCount   *int   `json:"watcher_count,omitempty"`
	PlayCount      *int   `json:"play_count,omitempty"`
	CollectedCount *int   `json:"collected_count,omitempty"`
	ListCount      *int   `json:"list_count,omitempty"`
	Movie          *Movie `json:"movie,omitempty"`
}

func (m MoviesItem) String() string {
	return Stringify(m)
}
