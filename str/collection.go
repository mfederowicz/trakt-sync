// Package str used for structs
package str

// Collection represents JSON collection object
type Collection struct {
	ItemCount *int `json:"item_count,omitempty"`
}

func (c Collection) String() string {
	return Stringify(c)
}
