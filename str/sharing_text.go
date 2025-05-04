// Package str used for structs
package str

// SharingText represents JSON sharing text object
type SharingText struct {
	Watching *string `json:"watching,omitempty"`
	Watched  *string `json:"watched,omitempty"`
	Rated    *string `json:"rated,omitempty"`
}

func (s SharingText) String() string {
	return Stringify(s)
}
