// Package str used for structs
package str

// Language represents JSON code object
type Language struct {
	Name *string `json:"name,omitempty"`
	Code *string `json:"code,omitempty"`
}

func (l Language) String() string {
	return Stringify(l)
}
