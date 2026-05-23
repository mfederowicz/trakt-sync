// Package str used for structs
package str

// PersonalListItem represents JSON list item object
type PersonalListItem struct {
	Notes *string `json:"notes,omitempty"`
}

func (p PersonalListItem) String() string {
	return Stringify(p)
}
