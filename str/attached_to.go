// Package str used for structs
package str

// AttachedTo represents JSON notes object
type AttachedTo struct {
	ID   *int    `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
}

func (a AttachedTo) String() string {
	return Stringify(a)
}
