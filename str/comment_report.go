// Package str used for structs
package str

// CommentReport represents JSON comment report object
type CommentReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (c CommentReport) String() string {
	return Stringify(c)
}
