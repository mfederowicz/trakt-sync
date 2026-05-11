// Package str used for structs
package str

// Comments represents JSON sesons object
type Comments struct {
	LikedAt   *Timestamp `json:"liked_at,omitempty"`
	BlockedAt *Timestamp `json:"blocked_at,omitempty"`
}

func (c Comments) String() string {
	return Stringify(c)
}
