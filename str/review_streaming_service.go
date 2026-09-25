// Package str used for structs
package str

// ReviewStreamingService represents JSON streaming service with per-type watched counts
type ReviewStreamingService struct {
	Source *string `json:"source,omitempty"`
	Name   *string `json:"name,omitempty"`
	Shows  *int    `json:"shows,omitempty"`
	Movies *int    `json:"movies,omitempty"`
	All    *int    `json:"all,omitempty"`
}

func (r ReviewStreamingService) String() string {
	return Stringify(r)
}
