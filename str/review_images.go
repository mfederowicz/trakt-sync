// Package str used for structs
package str

// ReviewImages represents JSON review images object
type ReviewImages struct {
	Cover *string `json:"cover,omitempty"`
	Story *string `json:"story,omitempty"`
}

func (r ReviewImages) String() string {
	return Stringify(r)
}
