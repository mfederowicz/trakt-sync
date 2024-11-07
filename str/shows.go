// Package str used for structs
package str

// Shows represents JSON shows object
type Shows struct {
	Watched   *int `json:"watched,omitempty"`
	Collected *int `json:"collected,omitempty"`
	Ratings   *int `json:"ratings,omitempty"`
	Comments  *int `json:"comments,omitempty"`
}

func (s Shows) String() string {
	return Stringify(s)
}
