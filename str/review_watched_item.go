// Package str used for structs
package str

// ReviewWatchedItem represents JSON first / last watched object: a movie or an episode with its show
type ReviewWatchedItem struct {
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Movie     *Movie     `json:"movie,omitempty"`
	Show      *Show      `json:"show,omitempty"`
	Episode   *Episode   `json:"episode,omitempty"`
}

func (r ReviewWatchedItem) String() string {
	return Stringify(r)
}
