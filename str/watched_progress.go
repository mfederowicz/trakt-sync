// Package str used for structs
package str

// WatchedProgress represents JSON show_collection_progress object
type WatchedProgress struct {
	Aired         *int       `json:"aired,omitempty"`
	Completed     *int       `json:"completed,omitempty"`
	LastWatchedAt *Timestamp `json:"last_watched_at,omitempty"`
	ResetAt       *Timestamp `json:"reset_at,omitempty"`
	Seasons       []*Season  `json:"seasons,omitempty"`
	HiddenSeasons []*Season  `json:"hidden_seasons,omitempty"`
	NextEpisode   *Episode   `json:"next_episode,omitempty"`
	LastEpisode   *Episode   `json:"last_episode,omitempty"`
}

func (w WatchedProgress) String() string {
	return Stringify(w)
}
