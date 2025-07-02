// Package str used for structs
package str

// Shows represents JSON shows object
type Shows struct {
	Watched       *int       `json:"watched,omitempty"`
	Collected     *int       `json:"collected,omitempty"`
	Ratings       *int       `json:"ratings,omitempty"`
	Comments      *int       `json:"comments,omitempty"`
	RatedAt       *Timestamp `json:"rated_at,omitempty"`
	WatchlistedAt *Timestamp `json:"watchlisted_at,omitempty"`
	FavoritedAt   *Timestamp `json:"favorited_at,omitempty"`
	CommentedAt   *Timestamp `json:"commented_at,omitempty"`
	HiddenAt      *Timestamp `json:"hidden_at,omitempty"`
	DroppedAt     *Timestamp `json:"dropped_at,omitempty"`
}

func (s Shows) String() string {
	return Stringify(s)
}
