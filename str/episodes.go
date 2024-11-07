// Package str used for structs
package str

// Episodes represents JSON episodes object
type Episodes struct {
	Plays     *int `json:"plays,omitempty"`
	Watched   *int `json:"watched,omitempty"`
	Minutes   *int `json:"minutes,omitempty"`
	Collected *int `json:"collected,omitempty"`
	Ratings   *int `json:"ratings,omitempty"`
	Comments  *int `json:"comments,omitempty"`
}

func (e Episodes) String() string {
	return Stringify(e)
}
