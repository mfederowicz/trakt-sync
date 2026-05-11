// Package str used for structs
package str

// Episodes represents JSON episodes object
type Episodes struct {
	Plays         *int       `json:"plays,omitempty"`
	Watched       *int       `json:"watched,omitempty"`
	Minutes       *int       `json:"minutes,omitempty"`
	Collected     *int       `json:"collected,omitempty"`
	Ratings       *int       `json:"ratings,omitempty"`
	Comments      *int       `json:"comments,omitempty"`
	WatchedAt     *Timestamp `json:"watched_at,omitempty"`
	CollectedAt   *Timestamp `json:"collected_at,omitempty"`
	RatedAt       *Timestamp `json:"rated_at,omitempty"`
	WatchlistedAt *Timestamp `json:"watchlisted_at,omitempty"`
	FavoritedAt   *Timestamp `json:"favorited_at,omitempty"`
	CommentedAt   *Timestamp `json:"commented_at,omitempty"`
	PausedAt      *Timestamp `json:"paused_at,omitempty"`
}

func (e Episodes) String() string {
	return Stringify(e)
}
