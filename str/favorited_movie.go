// Package str used for structs
package str

// FavoritedMovie represents JSON favorited movie object
type FavoritedMovie struct {
	UserCount *int   `json:"user_count,omitempty"`
	Movie     *Movie `json:"movie,omitempty"`
}

func (m FavoritedMovie) String() string {
	return Stringify(m)
}
