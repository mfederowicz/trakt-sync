// Package str used for structs
package str

// Notes represents JSON notes object
type Notes struct {
	ID        *int         `json:"id,omitempty"`
	Notes     *string      `json:"notes,omitempty"`
	Privacy   *string      `json:"privacy,omitempty"`
	ItemCount *int         `json:"item_count,omitempty"`
	Spoiler   *bool        `json:"spoiler,omitempty"`
	CreatedAt *Timestamp   `json:"created_at,omitempty"`
	UpdatedAt *Timestamp   `json:"updated_at,omitempty"`
	User      *UserProfile `json:"user,omitempty"`
	Movie     *Movie       `json:"movie,omitempty"`
	Show      *Show        `json:"show,omitempty"`
	Season    *Season      `json:"season,omitempty"`
	Episode   *Episode     `json:"episode,omitempty"`
	Person    *Person      `json:"person,omitempty"`
}

func (n Notes) String() string {
	return Stringify(n)
}
