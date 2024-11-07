// Package str used for structs
package str

// Ratings represents JSON ratings object
type Ratings struct {
	Total        *int            `json:"total,omitempty"`
	Distribution *map[string]int `json:"distribution,omitempty"`
}

func (r Ratings) String() string {
	return Stringify(r)
}
