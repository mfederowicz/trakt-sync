// Package str used for structs
package str

// TrendingMovie represents JSON trendig movie object
type TrendingMovie struct {
	Watchers *int   `json:"watchers,omitempty"`
	Movie    *Movie `json:"movie,omitempty"`
}

func (t TrendingMovie) String() string {
	return Stringify(t)
}
