// Package str used for structs
package str

// ReviewStreamingServices represents JSON subscription services the user watched on (month in review only)
type ReviewStreamingServices struct {
	Country  *string                   `json:"country,omitempty"`
	Services []*ReviewStreamingService `json:"services,omitempty"`
}

func (r ReviewStreamingServices) String() string {
	return Stringify(r)
}
