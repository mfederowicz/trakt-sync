// Package str used for structs
package str

// ListComment represents JSON list comment object
type ListComment struct {
	ID        *int         `json:"id,omitempty"`
	ParentID  *int         `json:"parent_id,omitempty"`
	CreatedAt *Timestamp   `json:"created_at,omitempty"`
	UpdatedAt *Timestamp   `json:"updated_at,omitempty"`
	Comment   *string      `json:"comment,omitempty"`
	Spoiler   *bool        `json:"spoiler,omitempty"`
	Review    *bool        `json:"review,omitempty"`
	Replies   *int         `json:"replies,omitempty"`
	Likes     *int         `json:"likes,omitempty"`
	UserStats *UserStats   `json:"user_stats,omitempty"`
	User      *UserProfile `json:"user,omitempty"`
}

func (l ListComment) String() string {
	return Stringify(l)
}
