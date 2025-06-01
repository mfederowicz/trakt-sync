// Package str used for structs
package str

// Lists represents JSON lists object
type Lists struct {
	LikedAt     *Timestamp `json:"liked_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`
	CommentedAt *Timestamp `json:"commented_at,omitempty"`
}

func (l Lists) String() string {
	return Stringify(l)
}
