// Package str used for structs
package str

// OutputEpisode represents JSON episode object used in deduplication
type OutputEpisode struct {
	Title     *string    `json:"title,omitempty"`
	Rating    *int       `json:"rating,omitempty"`
	Year      *int       `json:"year,omitempty"`
	IDs       *IDs       `json:"ids,omitempty"`
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
	RatedAt   *Timestamp `json:"rated_at,omitempty"`
}

func (o OutputEpisode) String() string {
	return Stringify(o)
}
