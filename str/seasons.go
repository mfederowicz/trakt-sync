// Package str used for structs
package str

// Seasons represents JSON sesons object
type Seasons struct {
	Ratings       *int       `json:"ratings,omitempty"`
	Comments      *int       `json:"comments,omitempty"`
	RatedAt       *Timestamp `json:"rated_at,omitempty"`
	WatchlistedAt *Timestamp `json:"watchlisted_at,omitempty"`
	CommentedAt   *Timestamp `json:"commented_at,omitempty"`
	HiddenAt      *Timestamp `json:"hidden_at,omitempty"`
}

func (s Seasons) String() string {
	return Stringify(s)
}
