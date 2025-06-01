// Package str used for structs
package str

// EpisodeStats represents JSON episode stats object
type EpisodeStats struct {
	Watchers   *int `json:"watchers,omitempty"`
	Plays      *int `json:"plays,omitempty"`
	Collectors *int `json:"collectors,omitempty"`
	Comments   *int `json:"comments,omitempty"`
	Lists      *int `json:"lists,omitempty"`
	Votes      *int `json:"votes,omitempty"`
	Favorited  *int `json:"favorited,omitempty"`
}

func (s EpisodeStats) String() string {
	return Stringify(s)
}
