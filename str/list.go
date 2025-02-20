// Package str used for structs
package str

// List represents JSON list object
type List struct {
	LikeCount    *int          `json:"like_count,omitempty"`
	CommentCount *int          `json:"comment_count,omitempty"`
	List         *PersonalList `json:"list,omitempty"`
}

func (l List) String() string {
	return Stringify(l)
}
