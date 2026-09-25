// Package str used for structs
package str

// ReviewCountryCount represents JSON country counts object, sorted by count
type ReviewCountryCount struct {
	CountryCount *int             `json:"country_count,omitempty"`
	Countries    []*ReviewCountry `json:"countries,omitempty"`
}

func (r ReviewCountryCount) String() string {
	return Stringify(r)
}
