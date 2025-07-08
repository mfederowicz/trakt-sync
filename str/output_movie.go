// Package str used for structs
package str

// OutputMovie represents JSON movie object used in deduplication
type OutputMovie struct {
	Title     *string    `json:"title,omitempty"`
	Year      *int       `json:"year,omitempty"`
	IDs       *IDs       `json:"ids,omitempty"`
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
}

func (o OutputMovie) String() string {
	return Stringify(o)
}
