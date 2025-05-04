// Package str used for structs
package str

// UserListItem represents JSON user personal list item object
type UserListItem struct {
	Rank     *int       `json:"rank,omitempty"`
	ID       *int       `json:"id,omitempty"`
	ListedAt *Timestamp `json:"listed_at,omitempty"`
	Notes    *string    `json:"notes,omitempty"`
	Type     *string    `json:"type,omitempty"`
	Movie    *Movie     `json:"movie,omitempty"`
	Show     *Show      `json:"show,omitempty"`
	Season   *Season    `json:"season,omitempty"`
	Episode  *Episode   `json:"episode,omitempty"`
	Person   *Person    `json:"person,omitempty"`
}

func (i UserListItem) String() string {
	return Stringify(i)
}
