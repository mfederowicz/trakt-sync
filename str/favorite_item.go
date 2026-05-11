// Package str used for structs
package str

// FavoriteItem represents JSON favorite item object
type FavoriteItem struct {
	Notes *string `json:"notes,omitempty"`
}

func (f FavoriteItem) String() string {
	return Stringify(f)
}
