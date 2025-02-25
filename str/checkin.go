// Package str used for structs
package str

// CheckIn represents JSON checkin object
type CheckIn struct {
	ID      *int64     `json:"id,omitempty"`
	Watched *Timestamp `json:"watched,omitempty"`
	Expires *Timestamp `json:"expires_at,omitempty"`
	Movie   *Movie     `json:"movie,omitempty"`
	Show    *Show      `json:"show,omitempty"`
	Episode *Episode   `json:"episode,omitempty"`
	Sharing *Sharing   `json:"sharing,omitempty"`
	Message *string    `json:"message,omitempty"`
}

func (c CheckIn) String() string {
	return Stringify(c)
}
