// Package str used for structs
package str

// NotesItem represents JSON notes attached item object
type NotesItem struct {
	AttachedTo *AttachedTo `json:"attached_to,omitempty"`
	Type       *string     `json:"type,omitempty"`
	Movie      *Movie      `json:"movie,omitempty"`
	Show       *Show       `json:"show,omitempty"`
	Season     *Season     `json:"season,omitempty"`
	Episode    *Episode    `json:"episode,omitempty"`
	Person     *Person     `json:"person,omitempty"`
}

func (n NotesItem) String() string {
	return Stringify(n)
}
