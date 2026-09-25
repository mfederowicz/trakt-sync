// Package str used for structs
package str

// SmartListImages represents JSON smart list images object
type SmartListImages struct {
	Posters []string `json:"posters,omitempty"`
}

func (s SmartListImages) String() string {
	return Stringify(s)
}
