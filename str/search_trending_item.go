// Package str used for structs
package str

// SearchTrendingItem represents JSON trending search response object
type SearchTrendingItem struct {
	ID     *int64  `json:"id,omitempty"`
	Count  *int    `json:"count,omitempty"`
	Type   *string `json:"type,omitempty"`
	Movie  *Movie  `json:"movie,omitempty"`
	Show   *Show   `json:"show,omitempty"`
	Person *Person `json:"person,omitempty"`
}

func (i SearchTrendingItem) String() string {
	return Stringify(i)
}
