// Package str used for structs
package str

// HiddenItem represents JSON hidden item object
type HiddenItem struct {
	HiddenAt *Timestamp `json:"hidden_at,omitempty"`
	Type     *string    `json:"type,omitempty"`
	Movie    *Movie     `json:"movie,omitempty"`
	Season   *Season    `json:"season,omitempty"`
	Episode  *Episode   `json:"episode,omitempty"`
	Show     *Show      `json:"show,omitempty"`
	List     *List      `json:"list,omitempty"`
	Comment  *Comment   `json:"comment,omitempty"`
}

func (h HiddenItem) String() string {
	return Stringify(h)
}
