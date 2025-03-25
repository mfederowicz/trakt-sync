// Package str used for structs
package str

// Studio represents JSON studio object
type Studio struct {
	Name    *string `json:"name,omitempty"`
	Country *string `json:"country,omitempty"`
	IDs     *IDs    `json:"ids,omitempty"`
}

func (s Studio) String() string {
	return Stringify(s)
}
