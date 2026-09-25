// Package str used for structs
package str

// ReviewCountries represents JSON watched shows and movies grouped by production country
type ReviewCountries struct {
	Shows  *ReviewCountryCount `json:"shows,omitempty"`
	Movies *ReviewCountryCount `json:"movies,omitempty"`
}

func (r ReviewCountries) String() string {
	return Stringify(r)
}
