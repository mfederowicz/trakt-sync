// Package str used for structs
package str

// SeasonStats represents JSON season stats object
type SeasonStats struct {
	Watchers   *int `json:"watchers,omitempty"`
	Plays      *int `json:"plays,omitempty"`
	Collectors *int `json:"collectors,omitempty"`
	Comments   *int `json:"comments,omitempty"`
	Lists      *int `json:"lists,omitempty"`
	Votes      *int `json:"votes,omitempty"`
	Favorited  *int `json:"favorited,omitempty"`
}

func (s SeasonStats) String() string {
	return Stringify(s)
}
