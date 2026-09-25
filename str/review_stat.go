// Package str used for structs
package str

// ReviewStat represents JSON review stat object with its total and averages
type ReviewStat struct {
	Total   *float64 `json:"total,omitempty"`
	Yearly  *float64 `json:"yearly,omitempty"`
	Monthly *float64 `json:"monthly,omitempty"`
	Weekly  *float64 `json:"weekly,omitempty"`
	Daily   *float64 `json:"daily,omitempty"`
}

func (r ReviewStat) String() string {
	return Stringify(r)
}
