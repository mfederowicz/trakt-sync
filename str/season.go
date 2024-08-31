// Package str used for structs
package str

// Season represents JSON season object
type Season struct {
	Number *int `json:"number,omitempty"`
	IDs    *IDs `json:"ids,omitempty"`
}

func (s Season) String() string {
	return Stringify(s)
}
