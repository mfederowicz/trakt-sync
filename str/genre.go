// Package str used for structs
package str

// Genre represents JSON slug object
type Genre struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

func (g Genre) String() string {
	return Stringify(g)
}
