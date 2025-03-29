// Package str used for structs
package str

// TvNetwork represents JSON tv network object
type TvNetwork struct {
	Name    *string `json:"name,omitempty"`
	Country *string `json:"country,omitempty"`
	IDs     *IDs    `json:"ids,omitempty"`
}

func (t TvNetwork) String() string {
	return Stringify(t)
}
