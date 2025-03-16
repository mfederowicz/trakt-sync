// Package str used for structs
package str

// PlayedMovie represents JSON played movie object
type PlayedMovie struct {
	WatcherCount   *int   `json:"watcher_count,omitempty"`
	PlayCount      *int   `json:"play_count,omitempty"`
	CollectedCount *int   `json:"collected_count,omitempty"`
	Movie          *Movie `json:"movie,omitempty"`
}

func (m PlayedMovie) String() string {
	return Stringify(m)
}
