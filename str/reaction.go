// Package str used for structs
package str

// Reaction represents JSON reaction object
type Reaction struct {
	Type *string `json:"type,omitempty"`
}

func (r Reaction) String() string {
	return Stringify(r)
}
