// Package str used for structs
package str

// Seasons represents JSON sesons object
type Seasons struct {
	Ratings   *int `json:"ratings,omitempty"`
	Comments  *int `json:"comments,omitempty"`
}

func (s Seasons) String() string {
	return Stringify(s)
}
