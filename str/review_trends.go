// Package str used for structs
package str

// ReviewTrends represents JSON most-watched shows and movies premiering each month
type ReviewTrends struct {
	Shows  []*ReviewTrend `json:"shows,omitempty"`
	Movies []*ReviewTrend `json:"movies,omitempty"`
}

func (r ReviewTrends) String() string {
	return Stringify(r)
}
