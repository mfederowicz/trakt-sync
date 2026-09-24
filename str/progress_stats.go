// Package str used for structs
package str

// ProgressStats represents JSON watch stats in show progress
type ProgressStats struct {
	PlayCount      *int `json:"play_count,omitempty"`
	MinutesWatched *int `json:"minutes_watched,omitempty"`
	MinutesLeft    *int `json:"minutes_left,omitempty"`
}

func (p ProgressStats) String() string {
	return Stringify(p)
}
