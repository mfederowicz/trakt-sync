// Package str used for structs
package str

// ReviewStats represents JSON review stats object for all media, shows and movies
type ReviewStats struct {
	All    *ReviewStatsCategories `json:"all,omitempty"`
	Shows  *ReviewStatsCategories `json:"shows,omitempty"`
	Movies *ReviewStatsCategories `json:"movies,omitempty"`
}

func (r ReviewStats) String() string {
	return Stringify(r)
}
