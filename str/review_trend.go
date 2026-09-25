// Package str used for structs
package str

// ReviewTrend represents JSON trend object; watchers is the global play count, watched tells if this user watched it
type ReviewTrend struct {
	Month     *int    `json:"month,omitempty"`
	MonthName *string `json:"month_name,omitempty"`
	Watchers  *int    `json:"watchers,omitempty"`
	Watched   *bool   `json:"watched,omitempty"`
	Show      *Show   `json:"show,omitempty"`
	Movie     *Movie  `json:"movie,omitempty"`
}

func (r ReviewTrend) String() string {
	return Stringify(r)
}
