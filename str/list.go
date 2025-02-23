// Package str used for structs
package str

// List represents JSON list object
type List struct {
	Count        *int          `json:"count,omitempty"`
	ItemCount    *int          `json:"item_count,omitempty"`
	LikeCount    *int          `json:"like_count,omitempty"`
	CommentCount *int          `json:"comment_count,omitempty"`
	List         *PersonalList `json:"list,omitempty"`
}

func (l List) String() string {
	return Stringify(l)
}
