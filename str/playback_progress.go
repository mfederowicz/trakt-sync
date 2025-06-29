// Package str used for structs
package str

// PlaybackProgress represents JSON playback_progress object
type PlaybackProgress struct {
	Progress *float32   `json:"aired,omitempty"`
	PausedAt *Timestamp `json:"paused_at,omitempty"`
	ID       *int64     `json:"id,omitempty"`
	Type     *string    `json:"type,omitempty"`
	Episode  *Episode   `json:"episode,omitempty"`
	Show     *Show      `json:"show,omitempty"`
	Movie    *Movie     `json:"movie,omitempty"`
}

func (s PlaybackProgress) String() string {
	return Stringify(s)
}
