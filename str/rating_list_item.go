// Package str used for structs
package str

// RatingListItem represents JSON list object
type RatingListItem struct {
	ID      *int64     `json:"id,omitempty"`
	RatedAt *Timestamp `json:"rated_at,omitempty"`
	Rating  *int       `json:"rating,omitempty"`
	Type    *string    `json:"type,omitempty"`
	Movie   *Movie     `json:"movie,omitempty"`
	Episode *Episode   `json:"episode,omitempty"`
	Show    *Show      `json:"show,omitempty"`
	Season  *Season    `json:"season,omitempty"`
}

func (r RatingListItem) String() string {
	return Stringify(r)
}
