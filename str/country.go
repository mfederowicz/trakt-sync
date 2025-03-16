// Package str used for structs
package str

// Country represents JSON country object
type Country struct {
	Name *string `json:"name,omitempty"`
	Code *string `json:"code,omitempty"`
}

func (c Country) String() string {
	return Stringify(c)
}
