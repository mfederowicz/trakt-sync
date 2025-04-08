// Package str used for structs
package str

// Scrobble represents JSON scrobble object
type Scrobble struct {
	ID       *int     `json:"id,omitempty"`
	Action   *string  `json:"action,omitempty"`
	Progress *float64 `json:"progress,omitempty"`
	Sharing  *Sharing `json:"sharing,omitempty"`
	Movie    *Movie   `json:"movie,omitempty"`
	Episode  *Episode `json:"episode,omitempty"`
	Show     *Show    `json:"show,omitempty"`
}

func (s Scrobble) String() string {
	return Stringify(s)
}
