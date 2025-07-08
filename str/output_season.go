// Package str used for structs
package str

// OutputSeason represents JSON season object used in deduplication
type OutputSeason struct {
	Title     *string    `json:"title,omitempty"`
	Year      *int       `json:"year,omitempty"`
	IDs       *IDs       `json:"ids,omitempty"`
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
}

func (o OutputSeason) String() string {
	return Stringify(o)
}
