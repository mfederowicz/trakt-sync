// Package str used for structs
package str

// Comment represents JSON comment object
type Comment struct {
	ID        *int          `json:"id,omitempty"`
	ParentID  *int          `json:"parent_id,omitempty"`
	CreatedAt *Timestamp    `json:"created_at,omitempty"`
	UpdatedAt *Timestamp    `json:"updated_at,omitempty"`
	Comment   *string       `json:"comment,omitempty"`
	Spoiler   *bool         `json:"spoiler,omitempty"`
	Sharing   *Sharing      `json:"sharing,omitempty"`
	Review    *bool         `json:"review,omitempty"`
	Replies   *int          `json:"replies,omitempty"`
	Likes     *int          `json:"likes,omitempty"`
	UserStats *UserStats    `json:"user_stats,omitempty"`
	User      *UserProfile  `json:"user,omitempty"`
	Movie     *Movie        `json:"movie,omitempty"`
	Show      *Show         `json:"show,omitempty"`
	Season    *Season       `json:"season,omitempty"`
	Episode   *Episode      `json:"episode,omitempty"`
	List      *PersonalList `json:"list,omitempty"`
}

func (c Comment) String() string {
	return Stringify(c)
}
