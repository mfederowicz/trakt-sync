// Package str used for structs
package str

// WatchNowRank represents JSON watch now streaming rank object
type WatchNowRank struct {
	Rank  *int    `json:"rank,omitempty"`
	Delta *int    `json:"delta,omitempty"`
	Link  *string `json:"link,omitempty"`
}

func (w WatchNowRank) String() string {
	return Stringify(w)
}
