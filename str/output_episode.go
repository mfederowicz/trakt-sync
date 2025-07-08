// Package str used for structs
package str

// Out represents JSON season object used in deduplication
type OutputEpisode struct {
	Title     *string    `json:"title,omitempty"`
	Year      *int       `json:"year,omitempty"`
	IDs       *IDs       `json:"ids,omitempty"`
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
}

func (o OutputEpisode) String() string {
	return Stringify(o)
}
