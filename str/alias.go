// Package str used for structs
package str

// Alias represents JSON alias object
type Alias struct {
	Title   *string `json:"title,omitempty"`
	Country *string `json:"country,omitempty"`
}

func (a Alias) String() string {
	return Stringify(a)
}
