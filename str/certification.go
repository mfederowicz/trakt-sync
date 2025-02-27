// Package str used for structs
package str

// Certification represents JSON certification object
type Certification struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (c Certification) String() string {
	return Stringify(c)
}
