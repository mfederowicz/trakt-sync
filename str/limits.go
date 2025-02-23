// Package str used for structs
package str

// Limits represents JSON limits object
type Limits struct {
	List       *List       `json:"list,omitempty"`
	Watchlist  *Watchlist  `json:"watchlist,omitempty"`
	Favorites  *Favorites  `json:"favorites,omitempty"`
	Search     *Search     `json:"search,omitempty"`
	Collection *Collection `json:"collection,omitempty"`
	Notes      *Notes      `json:"notes,omitempty"`
}

func (l Limits) String() string {
	return Stringify(l)
}
