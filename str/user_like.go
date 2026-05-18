// Package str used for structs
package str

// UserLike represents JSON user like object
type UserLike struct {
	LikedAt     *Timestamp    `json:"liked_at,omitempty"`
	Type        *string       `json:"type,omitempty"`
	List        *PersonalList `json:"list,omitempty"`
	CommentType *string       `json:"comment_type,omitempty"`
	Comment     *Comment      `json:"comment,omitempty"`
	Movie       *Movie        `json:"movie,omitempty"`
	Season      *Season       `json:"season,omitempty"`
	Show        *Show         `json:"show,omitempty"`
}

func (u UserLike) String() string {
	return Stringify(u)
}
