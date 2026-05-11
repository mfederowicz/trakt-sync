// Package str used for structs
package str

// Movies represents JSON movies object
type Movies struct {
	Watched       *int       `json:"watched,omitempty"`
	Collected     *int       `json:"collected,omitempty"`
	Plays         *int       `json:"plays,omitempty"`
	Minutes       *int       `json:"minutes,omitempty"`
	Ratings       *int       `json:"ratings,omitempty"`
	Comments      *int       `json:"comments,omitempty"`
	WatchedAt     *Timestamp `json:"watched_at,omitempty"`
	CollectedAt   *Timestamp `json:"collected_at,omitempty"`
	RatedAt       *Timestamp `json:"rated_at,omitempty"`
	WatchlistedAt *Timestamp `json:"watchlisted_at,omitempty"`
	FavoritedAt   *Timestamp `json:"favorited_at,omitempty"`
	CommentedAt   *Timestamp `json:"commented_at,omitempty"`
	PausedAt      *Timestamp `json:"paused_at,omitempty"`
	HiddenAt      *Timestamp `json:"hidden_at,omitempty"`
}

func (m Movies) String() string {
	return Stringify(m)
}
