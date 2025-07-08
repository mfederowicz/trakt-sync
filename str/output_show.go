// Package str used for structs
package str

// OutputShow represents JSON show object used in deduplication
type OutputShow struct {
	Title     *string    `json:"title,omitempty"`
	Year      *int       `json:"year,omitempty"`
	IDs       *IDs       `json:"ids,omitempty"`
	Seasons   *[]Season  `json:"seasons,omitempty"`
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
}

func (o OutputShow) String() string {
	return Stringify(o)
}
