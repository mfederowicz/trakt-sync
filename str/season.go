// Package str used for structs
package str

// Season represents JSON season object
type Season struct {
	Number   *int       `json:"number,omitempty"`
	Episodes []*Episode `json:"episodes,omitempty"`
	IDs      *IDs       `json:"ids,omitempty"`
}

func (s Season) String() string {
	return Stringify(s)
}
