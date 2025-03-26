// Package str used for structs
package str

// MovieStats represents JSON movie stats object
type MovieStats struct {
	Watchers   *int `json:"watchers,omitempty"`
	Plays      *int `json:"plays,omitempty"`
	Collectors *int `json:"collectors,omitempty"`
	Comments   *int `json:"comments,omitempty"`
	Lists      *int `json:"lists,omitempty"`
	Votes      *int `json:"votes,omitempty"`
	Favorited  *int `json:"favorited,omitempty"`
}

func (m MovieStats) String() string {
	return Stringify(m)
}
