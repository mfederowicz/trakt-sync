// Package str used for structs
package str

// Watchlist represents JSON watchlist object
type Watchlist struct {
	ItemCount *int `json:"item_count,omitempty"`
}

func (w Watchlist) String() string {
	return Stringify(w)
}
