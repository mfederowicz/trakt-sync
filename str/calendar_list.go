// Package str used for structs
package str

// CalendarList represents JSON calendar list item
type CalendarList struct {
	Released   *string    `json:"released,omitempty"`
	FirstAired *Timestamp `json:"first_aired,omitempty"`
	Episode    *Episode   `json:"episode,omitempty"`
	Show       *Show      `json:"show,omitempty"`
	Movie      *Movie     `json:"movie,omitempty"`
}

func (c CalendarList) String() string {
	return Stringify(c)
}
