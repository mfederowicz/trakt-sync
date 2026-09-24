// Package str used for structs
package str

// WatchNowPrices represents JSON watch now prices object
type WatchNowPrices struct {
	Rent     *string `json:"rent,omitempty"`
	Purchase *string `json:"purchase,omitempty"`
}

func (w WatchNowPrices) String() string {
	return Stringify(w)
}
