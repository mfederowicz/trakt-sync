// Package str used for structs
package str

// Movies represents JSON movies object
type Movies struct {
	Plays     *int `json:"plays,omitempty"`
	Watched   *int `json:"watched,omitempty"`
	Minutes   *int `json:"minutes,omitempty"`
	Collected *int `json:"collected,omitempty"`
	Ratings   *int `json:"ratings,omitempty"`
	Comments  *int `json:"comments,omitempty"`
}

func (m Movies) String() string {
	return Stringify(m)
}
