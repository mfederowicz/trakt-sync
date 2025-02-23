// Package str used for structs
package str

// Notes represents JSON notes object
type Notes struct {
	ItemCount *int `json:"item_count,omitempty"`
}

func (n Notes) String() string {
	return Stringify(n)
}
