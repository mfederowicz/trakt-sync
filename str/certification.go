// Package str used for structs
package str

// Certification represents JSON certification object
type Certification struct {
	Country       *string `json:"country,omitempty"`
	Certification *string `json:"certification,omitempty"`
	Name          *string `json:"name,omitempty"`
	Slug          *string `json:"slug,omitempty"`
	Description   *string `json:"description,omitempty"`
}

func (c Certification) String() string {
	return Stringify(c)
}
