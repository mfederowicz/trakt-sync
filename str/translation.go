// Package str used for structs
package str

// Translation represents JSON Translation object
type Translation struct {
	Title    *string `json:"title,omitempty"`
	Overview *string `json:"overview,omitempty"`
	Tagline  *string `json:"tagline,omitempty"`
	Language *string `json:"language,omitempty"`
	Country  *string `json:"country,omitempty"`
}

func (t Translation) String() string {
	return Stringify(t)
}
