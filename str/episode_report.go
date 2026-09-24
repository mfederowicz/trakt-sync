// Package str used for structs
package str

// EpisodeReport represents JSON episode report object
type EpisodeReport struct {
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (s EpisodeReport) String() string {
	return Stringify(s)
}
