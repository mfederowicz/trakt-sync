// Package str used for structs
package str

// ReviewCountry represents JSON country with its watched count
type ReviewCountry struct {
	Country *string `json:"country,omitempty"`
	Count   *int    `json:"count,omitempty"`
}

func (r ReviewCountry) String() string {
	return Stringify(r)
}
