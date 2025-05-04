// Package str used for structs
package str

// ShowsItem represents JSON movies item object
type ShowsItem struct {
	Watchers  *int  `json:"watchers,omitempty"`
	UserCount *int  `json:"user_count,omitempty"`
	Show      *Show `json:"show,omitempty"`
}

func (m ShowsItem) String() string {
	return Stringify(m)
}
