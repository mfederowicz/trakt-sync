// Package str used for structs
package str

// PersonalList represents JSON personal list object
type PersonalList struct {
	Name           *string      `json:"name,omitempty"`
	Description    *string      `json:"description,omitempty"`
	Privacy        *string      `json:"privacy,omitempty"`
	ShareLink      *string      `json:"share_link,omitempty"`
	Type           *string      `json:"type,omitempty"`
	DisplayNumbers *bool        `json:"display_strings,omitempty"`
	AllowComments  *bool        `json:"allow_comments,omitempty"`
	SortBy         *string      `json:"sort_by,omitempty"`
	SortHow        *string      `json:"sort_how,omitempty"`
	CreatedAt      *Timestamp   `json:"created_at,omitempty"`
	UpdatedAt      *Timestamp   `json:"updated_at,omitempty"`
	ItemCount      *int         `json:"item_count,omitempty"`
	CommentCount   *int         `json:"comment_count,omitempty"`
	Likes          *int         `json:"likes,omitempty"`
	IDs            *IDs         `json:"ids,omitempty"`
	User           *UserProfile `json:"user,omitempty"`
}

func (p PersonalList) String() string {
	return Stringify(p)
}
