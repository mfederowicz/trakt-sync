// Package str used for structs
package str

// WatchingResult represents JSON watching result object
type WatchingResult struct {
	ExpiresAt *Timestamp `json:"expires_at,omitempty"`
	StartedAt *Timestamp `json:"started_at,omitempty"`
	Action    *string    `json:"action,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Episode   *Episode   `json:"episode,omitempty"`
	Show      *Show      `json:"show,omitempty"`
	Movie     *Movie     `json:"movie,omitempty"`
}

func (w WatchingResult) String() string {
	return Stringify(w)
}
