// Package str used for structs
package str

// WatchlistItem represents JSON watchlist object
type WatchlistItem struct {
	Notes *string `json:"notes,omitempty"`
}

func (w WatchlistItem) String() string {
	return Stringify(w)
}
